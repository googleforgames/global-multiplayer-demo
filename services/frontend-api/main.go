// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"

	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	pballoc "agones.dev/agones/pkg/allocation/go"
	"github.com/pkg/errors"
	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc"
)

type OMServerResponse struct {
	IP     string
	Port   int
	Region string
}

type RegionalLatency struct {
	Region  string
	Latency int
}

/*
const (
	OMFrontendEndpoint = "open-match-frontend.cluster-name.svc.cluster.local:50504"
)
*/

/*

Dummy request until OpenMatch is implemented:

curl -X POST localhost:8080 \
     -H "Content-Type: application/json" \
     -d "[{\"Region\": \"EU\", \"Latency\": 50}, {\"Region\": \"US\", \"Latency\": 100}]"

Expected response:

{
	"IP":"127.0.0.1",
	"Port":7777,
	"Region":"EU"
}

*/

func main() {
	validateEnv()

	Port := os.Getenv("PORT")
	if Port == "" {
		Port = "8080"
		log.Printf("Defaulting to port %s\n", Port)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "{\"error\": \"%s\"}", err)
				log.Println("panic occurred:", err)
			}
		}()

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		LatencyList := parseRegionalLatencies(w, r)
		for _, LatencyItem := range LatencyList {
			log.Printf("Latency for %s is %d\n", LatencyItem.Region, LatencyItem.Latency)
		}

		// Sort by latency
		sort.SliceStable(LatencyList, func(i, j int) bool {
			return LatencyList[i].Latency < LatencyList[j].Latency
		})

		mOMResponse, _ := json.Marshal(findMatchingServer(LatencyList))

		fmt.Fprint(w, string(mOMResponse))
	})

	log.Printf("Listening on port %s\n", Port)
	http.ListenAndServe(":"+Port, nil)
}

func findMatchingServer(LatencyList []RegionalLatency) OMServerResponse {
	log.Printf("Looking for a server in the %s region.\n", LatencyList[0].Region)

	// TODO: Query OpenMatch on `OMFrontendEndpoint` in a specified region for a server. For now we query Agones directly.

	IP, Port := getAvailableServer()
	/* ip := "127.0.0.1"
	port := 7777 */

	return OMServerResponse{
		IP:     IP,
		Port:   Port,
		Region: LatencyList[0].Region}

}

func parseRegionalLatencies(rw http.ResponseWriter, req *http.Request) []RegionalLatency {
	decoder := json.NewDecoder(req.Body)
	var LatencyList []RegionalLatency
	err := decoder.Decode(&LatencyList)

	if err != nil {
		panic(err)
	}

	return LatencyList
}

func validateEnv() {
	keyFile := os.Getenv("key")
	certFile := os.Getenv("cert")
	cacertFile := os.Getenv("cacert")
	externalIP := os.Getenv("ip")

	if len(keyFile) == 0 || len(certFile) == 0 || len(cacertFile) == 0 || len(externalIP) == 0 {
		fmt.Println("key, cert, cacert and ip are mandatory ENV variables")
		os.Exit(0)
	}
}

func getAvailableServer() (string, int) {
	keyFile := os.Getenv("key")
	certFile := os.Getenv("cert")
	cacertFile := os.Getenv("cacert")
	externalIP := os.Getenv("ip")
	port := os.Getenv("port")
	namespace := os.Getenv("namespace")
	multicluster := os.Getenv("multicluster")

	if len(port) == 0 {
		port = "443"
	}

	if len(namespace) == 0 {
		namespace = "default"
	}

	if len(multicluster) == 0 {
		multicluster = "false"
	}

	endpoint := externalIP + ":" + port
	cert, err := ioutil.ReadFile(certFile)
	if err != nil {
		panic(err)
	}
	key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		panic(err)
	}
	cacert, err := ioutil.ReadFile(cacertFile)
	if err != nil {
		panic(err)
	}

	_multicluster, _ := strconv.ParseBool(multicluster)
	request := &pballoc.AllocationRequest{
		Namespace: namespace,
		MultiClusterSetting: &pballoc.MultiClusterSetting{
			Enabled: _multicluster,
		},
	}

	dialOpts, err := createRemoteClusterDialOption(cert, key, cacert)
	if err != nil {
		panic(err)
	}
	conn, err := grpc.Dial(endpoint, dialOpts)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	grpcClient := pballoc.NewAllocationServiceClient(conn)
	response, err := grpcClient.Allocate(context.Background(), request)
	if err != nil {
		fmt.Printf("%s", err)
		panic("Unable to allocate a server")
	}

	fmt.Printf("%s", response)

	return response.Address, int(response.Ports[len(response.Ports)-1].Port)
}

// createRemoteClusterDialOption creates a grpc client dial option with TLS configuration.
func createRemoteClusterDialOption(clientCert, clientKey, caCert []byte) (grpc.DialOption, error) {
	// Load client cert
	cert, err := tls.X509KeyPair(clientCert, clientKey)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
	if len(caCert) != 0 {
		// Load CA cert, if provided and trust the server certificate.
		// This is required for self-signed certs.
		tlsConfig.RootCAs = x509.NewCertPool()
		if !tlsConfig.RootCAs.AppendCertsFromPEM(caCert) {
			return nil, errors.New("only PEM format is accepted for server CA")
		}
	}

	return grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)), nil
}
