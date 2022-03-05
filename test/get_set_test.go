package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	pb "storage/kvstore"

	commom "storage/util/common"
	"storage/util/crc16"

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

	var statusCode commom.Status
	var errorCode commom.Error
	statusCode.Init()
	errorCode.Init()

	serverAddress := make(map[int64]string)
	serverStatus := make(map[int64]string)
	serverMaxKey := make(map[int64]int64)

	address := "192.168.1.128:50050"
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	check(err)
	defer conn.Close()
	c := pb.NewStorageClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Second))
	defer cancel()
	reply, err := c.GetConf(ctx, &pb.GetConfRequest{})
	check(err)
	if reply.GetStatus() != statusCode.Ok {
		fmt.Println("同步服务器配置出错...")
	}

	next := reply.GetNext()
	level := reply.GetLevel()
	hashSize := reply.GetHashSize()
	for _, server := range reply.GetServer() {
		idx := server.Idx
		serverAddress[idx] = server.Address
		serverMaxKey[idx] = server.MaxKey
		serverStatus[idx] = server.Status
	}

	hashFunc := func(key string) int64 {
		posCRC16 := int64(crc16.Checksum([]byte(key), crc16.IBMTable))
		pos := posCRC16 % (int64(math.Pow(2.0, float64(level))) * hashSize)
		if pos < next { // 分裂过了的
			pos = posCRC16 % (int64(math.Pow(2.0, float64(level+1))) * hashSize)
		}
		return pos
	}

	keys := []string{}
	values := []string{}

	path, err := os.Getwd()
	check(err)
	filename := filepath.Join(path, "../data/data2")
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		kv := strings.Split(sc.Text(), " ")
		keys = append(keys, kv[0])
		values = append(values, kv[1])
	}

	// SET
	for i := 0; i < len(keys); i++ {
		key := keys[i]
		value := values[i]
		idx := hashFunc(key)
		addr := serverAddress[idx]
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		check(err)
		defer conn.Close()
		c := pb.NewStorageClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Second))
		defer cancel()
		reply, err := c.Set(ctx, &pb.SetRequest{Key: key, Value: value})
		if err != nil || reply.GetStatus() != statusCode.Ok {
			panic(reply.GetErr())
		}
	}

	// GET
	for i := 0; i < len(keys); i++ {
		key := keys[i]
		value := values[i]
		idx := hashFunc(key)
		addr := serverAddress[idx]
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		check(err)
		defer conn.Close()
		c := pb.NewStorageClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Second))
		defer cancel()
		reply, err := c.Get(ctx, &pb.GetRequest{Key: key})
		if err != nil || reply.GetStatus() != statusCode.Ok {
			panic(reply.GetErr())
		}
		if reply.GetResult() != value {
			panic("Wrong value!!!")
		}
	}

	// DEL
	for i := 0; i < len(keys); i++ {
		key := keys[i]
		idx := hashFunc(key)
		addr := serverAddress[idx]
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		check(err)
		defer conn.Close()
		c := pb.NewStorageClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Second))
		defer cancel()
		reply, err := c.Del(ctx, &pb.DelRequest{Key: key})
		if err != nil || reply.GetStatus() != statusCode.Ok {
			panic(reply.GetErr())
		}
		if reply.GetResult() != 1 {
			log.Fatalf("Delete key %s failed!", key)
		}
	}

	// GET
	for i := 0; i < len(keys); i++ {
		key := keys[i]
		idx := hashFunc(key)
		addr := serverAddress[idx]
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		check(err)
		defer conn.Close()
		c := pb.NewStorageClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Second))
		defer cancel()
		reply, err := c.Get(ctx, &pb.GetRequest{Key: key})
		if err != nil {
			panic(err)
		}
		if reply.GetErr() != errorCode.NotFound {
			panic("Wrong value!!!")
		}
	}
}

func TestSplit(t *testing.T) {
	conn, err := grpc.Dial("192.168.1.128:"+fmt.Sprintf("%d", 50051), grpc.WithTransportCredentials(insecure.NewCredentials()))
	check(err)
	c := pb.NewStorageClient(conn)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Second))
	defer cancel()

	currentPath, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	filename := filepath.Join(currentPath, "data/benchmark_data")
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		operation := strings.Split(sc.Text(), " ")
		reply, _ := c.Get(ctx, &pb.GetRequest{Key: operation[1]})
		if reply.GetStatus() != "OK" || reply.GetResult() != operation[2] {
			t.Error("Wrong Answer!")
		}
	}
}
