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
	"strconv"
	"strings"
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
	currentPath, err := os.Getwd()
	check(err)
	filename := filepath.Join(currentPath, "data/benchmark_data")
	f, err := os.Open(filename)
	check(err)
	defer f.Close()

	sc := bufio.NewScanner(f)
	// data := make([][]string, 0)
	counter := 0

	// var wg sync.WaitGroup
	start := time.Now()
	for sc.Scan() {
		operation := strings.Split(sc.Text(), " ")
		// fmt.Println(operation)
		// index := crc16.Checksum([]byte(operation[0]), crc16.IBMTable) % 5
		// wg.Add(1)
		getServe(operation, conn)
		counter++
	}
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

func getServe(operation []string, conn *grpc.ClientConn) {
	// defer wg.Done()

	c := pb.NewStorageClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Second))
	defer cancel()

	opt := operation[0]
	switch opt {
	case "get":
		// continue
		reply, err := c.Get(ctx, &pb.GetRequest{Key: operation[1]})
		check(err)
		fmt.Println(reply.GetStatus())
	case "set":
		// continue
		reply, err := c.Set(ctx, &pb.SetRequest{Key: operation[1], Value: operation[2]})
		check(err)
		fmt.Println(reply.GetStatus())
		// fmt.Printf("%s\"\n", reply.GetStatus())
	case "del":
		// continue
		reply, err := c.Del(ctx, &pb.DelRequest{Key: operation[1]})
		check(err)
		fmt.Println(reply.GetStatus())
	case "split":
		port, err := strconv.Atoi(operation[1])
		reply, err := c.Split(ctx, &pb.SplitRequest{Port: int32(port)})
		check(err)
		fmt.Println(reply.GetResult())
	}

}
