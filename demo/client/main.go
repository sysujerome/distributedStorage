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
	"path"
	"strings"
	"time"

	pb "example.com/kvstore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	pwd  = flag.String("project_path", "/home/jerome/lab/distributedStorage/demo/", "the path of demo")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewStorageClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second))
	defer cancel()

	f, err := os.Open(path.Join(*pwd, "util/benchmark_data"))
	check(err)
	defer f.Close()
	sc := bufio.NewScanner(f)
	start := time.Now()
	counter := 0
	for sc.Scan() {
		operation := strings.Split(sc.Text(), " ")
		opt := operation[0]
		switch opt {
		case "get":
			_, err := c.Get(ctx, &pb.GetRequest{Key: operation[1]})
			check(err)
		case "set":
			_, err := c.Set(ctx, &pb.SetRequest{Key: operation[1], Value: operation[2]})
			check(err)
		case "del":
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
