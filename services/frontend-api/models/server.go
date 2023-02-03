// Copyright 2023 Google Inc. All Rights Reserved.
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

package models

import "log"

type OMServerResponse struct {
	IP   string
	Port int
}

type PingServer struct {
	IP     string
	Region string
}

func FindMatchingServer(regions []string) OMServerResponse {

	// TODO: Query OpenMatch on `OMFrontendEndpoint` in a preferred region for a server
	log.Printf("Looking for a server in the %s region.\n", regions[0])

	IP := "127.0.0.1"
	Port := 7777

	return OMServerResponse{
		IP:   IP,
		Port: Port}

}
