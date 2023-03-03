/*
 * Copyright 2023 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"github.com/gin-gonic/gin"
	"github.com/jellydator/ttlcache/v3"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
)

const (
	cacheKey        = "json"
	external        = "EXTERNAL"
	pingServiceName = "agones-ping-udp-service"
)

func main() {
	ctx := context.Background()

	forwardingClient, err := compute.NewForwardingRulesRESTClient(ctx)
	if err != nil {
		log.Fatalf("could not get forwarding client: %s", err)
	}
	defer forwardingClient.Close()
	cache := newCache(ctx, forwardingClient)

	r := gin.Default()
	if err := r.SetTrustedProxies(nil); err != nil {
		log.Fatalf("error setting trusted proxy: %s", err)
	}
	r.GET("/list", func(c *gin.Context) {
		services, err := getPingServicePerRegion(cache)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Printf("%v", services)
		c.JSON(http.StatusOK, services)
	})

	if err = r.Run(); err != nil {
		log.Fatal(err)
	}
}

// Service is the details of an individual Ping Service endpoint.
type Service struct {
	Name      string
	Namespace string
	Region    string
	Address   string
	Port      uint64
	Protocol  string
}

func newCache(ctx context.Context, forwardingClient *compute.ForwardingRulesClient) *ttlcache.Cache[string, []Service] {
	loader := ttlcache.LoaderFunc[string, []Service](
		func(c *ttlcache.Cache[string, []Service], key string) *ttlcache.Item[string, []Service] {

			services, err := listPingServices(ctx, forwardingClient)
			if err != nil {
				log.Printf("Error getting ping services: %s", err)
				return nil
			}

			// load from file/make an HTTP request
			item := c.Set(cacheKey, services, ttlcache.DefaultTTL)
			return item
		},
	)

	// Add some jitter to avoid multiple hits to the API at the same time on expiration.
	jitter := time.Duration(rand.Int63n(30)) * time.Second
	cache := ttlcache.New[string, []Service](
		ttlcache.WithLoader[string, []Service](loader),
		// only retrieve the information around once a minute
		ttlcache.WithTTL[string, []Service](time.Minute+jitter),
	)
	return cache
}

// getPingServicePerRegion returns one randomly selected UDP ping service per region for this project.
// The key in the map is the region, and the Service is the details of the Service.
func getPingServicePerRegion(c *ttlcache.Cache[string, []Service]) (map[string]Service, error) {
	item := c.Get(cacheKey)
	if item == nil {
		return nil, errors.New("cached service list not found")
	}

	list := item.Value()
	result := map[string]Service{}
	// randomise it so pings are routed to a variety of clusters.
	rand.Shuffle(len(list), func(i, j int) {
		list[i], list[j] = list[j], list[i]
	})

	for _, s := range list {
		if s.Name == pingServiceName {
			result[s.Region] = s
		}
	}

	return result, nil
}

func listPingServices(ctx context.Context, c *compute.ForwardingRulesClient) ([]Service, error) {
	list, err := listGKEServices(ctx, c)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, errors.New("no external endpoints found")
	}

	var result []Service
	for _, s := range list {
		if s.Name == pingServiceName {
			result = append(result, s)
		}
	}

	return result, nil
}

// listGKEServices lists all GKE Services that are exposed publicly in the project.
func listGKEServices(ctx context.Context, c *compute.ForwardingRulesClient) ([]Service, error) {
	project, err := findProject(ctx)
	if err != nil {
		return nil, err
	}

	req := &computepb.AggregatedListForwardingRulesRequest{
		Project: project,
	}

	var result []Service

	it := c.AggregatedList(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		for _, rule := range resp.Value.GetForwardingRules() {
			if rule.GetLoadBalancingScheme() == external {
				jsn := rule.GetDescription()
				if len(jsn) == 0 {
					continue
				}

				var desc map[string]string
				err := json.Unmarshal([]byte(jsn), &desc)
				if err != nil {
					log.Printf("Warning: Not JSON: %s", jsn)
				}

				serviceName, ok := desc["kubernetes.io/service-name"]
				if !ok {
					continue
				}
				split := strings.Split(serviceName, "/")
				if len(split) != 2 {
					continue
				}

				region := strings.Split(rule.GetRegion(), "/")

				svc := Service{
					Name:      split[1],
					Namespace: split[0],
					Region:    region[len(region)-1],
					Address:   rule.GetIPAddress(),
					Protocol:  rule.GetIPProtocol(),
				}
				// we know it's only one port in the range
				if ports := rule.GetPortRange(); len(ports) > 0 {
					// port range should be 5000-5000, so let's split and take the first one.
					ports := strings.Split(ports, "-")
					if len(ports) == 2 {
						svc.Port, err = strconv.ParseUint(ports[0], 10, 64)
						if err != nil {
							log.Printf("Warning: port not parseable: %s", ports[0])
						}
					} else {
						log.Printf("Warning: Port range format issue: %s", ports)
					}
				} else {
					log.Printf("Warning: No port range set for the UDP ping service")
				}

				log.Printf("Service: %+v", svc)

				result = append(result, svc)
			}
		}
	}

	return result, nil
}

// findProject uses several mechanisms to attempt to find the current project id.
func findProject(ctx context.Context) (string, error) {
	// Attempt to use the service accounts credentials.
	creds, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		return "", fmt.Errorf("could not find default credentials: %w", err)
	}

	// Local machine has no service credentials, so will return blank
	if len(creds.ProjectID) > 0 {
		log.Printf("Found project: [%s]", creds.ProjectID)
		return creds.ProjectID, nil
	}

	// Assume at this point we are on a local machine, with gcloud installed.
	cmd := exec.Command("gcloud", "info", "--format", "value(config.project)")
	var out strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error executing gcloud command (%s) to lookup project: %w. %s", cmd.String(), err, out.String())
	}

	project := strings.TrimSpace(out.String())
	log.Printf("Found project: [%s]", project)
	return project, nil
}
