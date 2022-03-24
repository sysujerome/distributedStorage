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
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	pb "example.com/kvstore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "192.168.1.128", "the address to connect to")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr+fmt.Sprintf(":%v", os.Args[1]), grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	check(err)
	// data := make([][]string, 0)
	counter := 0

	// var wg sync.WaitGroup
	start := time.Now()

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Second))
	defer cancel()
	c := pb.NewStorageClient(conn)
	reply, _ := c.Get(ctx, &pb.GetRequest{Key: "saa"})
	fmt.Printf("%v", reply.GetStatus())
	// log.Fatal(reply.GetResult())
	result, err := c.Scan(ctx, &pb.ScanRequest{Port: 50051})
	check(err)
	fmt.Printf("%v", result)
	duration := time.Since(start)
	fmt.Printf("dealing with %d operations took %v Seconds\n", counter, duration.Seconds())

	// wg.Wait()
	// for i := 0; i < 5; i++ {
	// 	conns[i].Close()
	// }
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
