package main

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	pb "example.com/kvstore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestSetGet(t *testing.T) {
	conns := make([]*grpc.ClientConn, 5)
	for i := 0; i < 5; i++ {
		var err error
		port := 50050 + i
		conns[i], err = grpc.Dial("192.168.1.128:"+fmt.Sprintf("%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		check(err)
	}
	for _, conn := range conns {
		c := pb.NewStorageClient(conn)
		// Contact the server and print out its response.
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Second))
		defer cancel()
		key1 := RandStringRunes(10)
		value1 := RandStringRunes(10)
		reply1, _ := c.Set(ctx, &pb.SetRequest{Key: key1, Value: value1})
		if reply1.GetStatus() != "OK" {
			t.Error(reply1.GetStatus())
		}
		reply2, _ := c.Get(ctx, &pb.GetRequest{Key: key1})
		// _, err := c.Get(ctx, &pb.GetRequest{Key: operation[1]})
		if reply2.GetStatus() != "OK" {
			t.Error(reply2.GetStatus())
		}
		if reply2.GetResult() != value1 {
			t.Error("Wrong Answer!!!")
		}
		reply3, _ := c.Del(ctx, &pb.DelRequest{Key: key1})
		if reply3.GetStatus() != "OK" {
			t.Error(reply2.GetStatus())
		}
		reply4, _ := c.Get(ctx, &pb.GetRequest{Key: key1})
		// _, err := c.Get(ctx, &pb.GetRequest{Key: operation[1]})
		if reply4.GetStatus() != "Not found!" {
			t.Error(reply4.GetStatus())
		}
		if reply4.GetResult() == value1 {
			t.Error("Wrong Answer!!!")
		}
	}
}
