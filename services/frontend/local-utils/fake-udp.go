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

// This is a local mock for getting regional ping servers as well as working UDP ping listener for client to connect to and measure latency.
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "%s - UDP echo server\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Usage: %s <IP> <portno>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "you have to specify an IP address, even if it's only 127.0.0.1\n")
		os.Exit(1)
	}

	ip := os.Args[1]
	port, _ := strconv.Atoi(os.Args[2])

	go handleHttp(ip, port)
	handleUdp(ip, port+1) // Just binding next port for UDP

}

func handleHttp(ip string, port int) {
	// http setup
	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		mock := `{
			"us-central1": {
				"Name": "agones-ping-udp-service",
				"Namespace": "agones-system",
				"Region": "my-local1",
				"Address": "127.0.0.1",
				"Port": ` + strconv.Itoa(port+1) + `,
				"Protocol": "UDP"
			},
			"europe-west1": {
				"Name": "agones-ping-udp-service",
				"Namespace": "agones-system",
				"Region": "my-local2",
				"Address": "localhost",
				"Port": ` + strconv.Itoa(port+1) + `,
				"Protocol": "UDP"
			}
		}`
		w.Write([]byte(mock))
	})

	log.Printf("Fake agones ping server running on: http://%s:%d/list", ip, port)
	http.ListenAndServe(ip+":"+strconv.Itoa(port), nil)
}

func handleUdp(ip string, port int) {
	// udp setup
	addr := net.UDPAddr{Port: port, IP: net.ParseIP(ip)}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatal(err)
	}

	b := make([]byte, 2048)

	log.Printf("Fake udp ping server running on: udp://%s:%d", ip, port)

	for {
		log.Printf("Accepting a new packet\n")

		cc, remote, rderr := conn.ReadFromUDP(b)

		if rderr != nil {
			log.Printf("net.ReadFromUDP() error: %s\n", rderr)
		} else {
			log.Printf("Read %d bytes from socket\n", cc)
			log.Printf("Bytes: %q\n", string(b[:cc]))
		}

		log.Printf("Remote address: %v\n", remote)

		cc, wrerr := conn.WriteTo(b[0:cc], remote)
		if wrerr != nil {
			log.Printf("net.WriteTo() error: %s\n", wrerr)
		} else {
			log.Printf("Wrote %d bytes to socket\n", cc)
		}
	}
}
