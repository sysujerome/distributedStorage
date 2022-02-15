/*
 *
 * Copyright 2015 gRPC authors.
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
 *
 */

// Package main implements a client for Greeter service.
package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	pb "example.com/kvstore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost", "the address to connect to")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	if len(os.Args) < 2 {
		log.Printf("usage: %s port number", string(os.Args[0]))
	}
	port := os.Args[1]
	if port == "50050" {
		fmt.Printf("store inside: ")
	} else {
		fmt.Printf("store in redis: ")
	}
	*addr = *addr + ":" + port
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewStorageClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Second))
	defer cancel()

	currentPath, err := os.Getwd()
	check(err)
	filename := filepath.Join(currentPath, "data/benchmark_data")
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	sc := bufio.NewScanner(f)
	// data := make([][]string, 0)
	counter := 0
	start := time.Now()
	for sc.Scan() {
		operation := strings.Split(sc.Text(), " ")
		// fmt.Println(operation)
		opt := operation[0]
		switch opt {
		case "get":
			// continue
			_, err := c.Get(ctx, &pb.GetRequest{Key: operation[1]})
			check(err)
		case "set":
			// continue
			_, err := c.Set(ctx, &pb.SetRequest{Key: operation[1], Value: operation[2]})
			check(err)
		case "del":
			// continue
			_, err := c.Del(ctx, &pb.DelRequest{Key: operation[1]})
			check(err)
		}
		counter++
	}
	duration := time.Since(start)
	fmt.Printf("dealing with %d operations took %v Seconds\n", counter, duration.Seconds())
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
