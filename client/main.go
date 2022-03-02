package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	pb "storage/kvstore"
	commom "storage/util/common"
	crc16 "storage/util/crc16"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	serverAddress map[int64]string
	serverStatus  map[int64]string
	serverMaxKey  map[int64]int64
	next          int64
	level         int64
	hashSize      int64
	ip            = flag.String("ip", "", "the base server to get configure")
	port          = flag.Int64("port", 0, "the base server to get configure")
	conf          = flag.String("conf", "", "the defaulted configure file")
	statusCode    commom.Status
	errorCode     commom.Error
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	flag.Parse()
	initConf()

	printConf()
	var conns []*grpc.ClientConn
	for _, addr := range serverAddress {
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		check(err)
		conns = append(conns, conn)
	}

	currentPath, err := os.Getwd()
	check(err)
	filename := filepath.Join(currentPath, "data/benchmark_data")
	f, err := os.Open(filename)
	check(err)
	defer f.Close()

	sc := bufio.NewScanner(f)
	counter := 0
	start := time.Now()
	for sc.Scan() {
		operation := strings.Split(sc.Text(), " ")
		key := operation[1]
		idx := hashFunc(key)
		getServe(operation, conns[idx])
		counter++
	}
	duration := time.Since(start)
	fmt.Printf("dealing with %d operations took %v Seconds\n", counter, duration.Seconds())
}

func getServe(operation []string, conn *grpc.ClientConn) {

	c := pb.NewStorageClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Second))
	defer cancel()

	opt := operation[0]
	switch opt {
	case "get":
		reply, err := c.Get(ctx, &pb.GetRequest{Key: operation[1]})
		check(err)
		fmt.Println(reply.GetStatus())
	case "set":
		reply, err := c.Set(ctx, &pb.SetRequest{Key: operation[1], Value: operation[2]})
		check(err)
		fmt.Println(reply.GetStatus())
	case "del":
		reply, err := c.Del(ctx, &pb.DelRequest{Key: operation[1]})
		check(err)
		fmt.Println(reply.GetStatus())
	}

}

func initConf() {

	serverAddress = make(map[int64]string)
	serverStatus = make(map[int64]string)
	serverMaxKey = make(map[int64]int64)
	address := ""
	defer flag.Parse()
	curPath, err := os.Getwd()
	check(err)
	if *conf != errorCode.NotDefined {
		filePath := filepath.Join(curPath, *conf)
		//检查文件是否存在
		if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
			panic("配置文件不存在...")
		}
		fileData, err := ioutil.ReadFile(filePath)
		check(err)
		var configure map[string]interface{}
		json.Unmarshal(fileData, &configure)
		shardConfs := configure["shard_confs"].([]interface{}) // shard_confs
		for _, v := range shardConfs {
			shardConf := v.(map[string]interface{})
			idx := int64(shardConf["shard_idx"].(float64))
			ShardNodeConfs := shardConf["shard_node_confs"].(map[string]interface{})
			for _, ShardNodeConf := range ShardNodeConfs {
				configureDetail := ShardNodeConf.(map[string]interface{})
				ip := configureDetail["ip"].(string)
				port := int64(configureDetail["base_port"].(float64))
				address := fmt.Sprintf("%s:%d", ip, port)
				serverAddress[idx] = address
				serverStatus[idx] = configureDetail["status"].(string)
				serverMaxKey[idx] = int64(configureDetail["max_key"].(float64))
			}
		}
		address = serverAddress[0]
	}
	if *ip != "" {
		address = fmt.Sprintf("%s:%d", *ip, *port)
	}
	if address == "" {
		fmt.Println("请输入ip和port, 或者输入配置文件路径")
	}
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
	next = reply.GetNext()
	level = reply.GetLevel()
	hashSize = reply.GetHashSize()
	for idx, server := range reply.GetServer() {
		serverAddress[int64(idx)] = server.Address
		serverMaxKey[int64(idx)] = server.MaxKey
		serverStatus[int64(idx)] = server.Status
	}
}

func printConf() {
	for idx, addr := range serverAddress {
		fmt.Printf("%d %s\t", idx, addr)
	}
	fmt.Println()

	for idx, maxKey := range serverMaxKey {
		fmt.Printf("%d %d\t", idx, maxKey)
	}
	fmt.Println()

	for idx, status := range serverStatus {
		fmt.Printf("%d %s\t", idx, status)
	}
	fmt.Println()

	fmt.Printf("next:\n%d\n", next)
	fmt.Printf("level:\n%d\n", level)
	fmt.Printf("hashSize:\n%d\n", hashSize)
	fmt.Printf("conf:\n%s\n", *conf)
}

func hashFunc(key string) int64 {
	posCRC16 := int64(crc16.Checksum([]byte(key), crc16.IBMTable))
	pos := posCRC16 % (int64(math.Pow(2.0, float64(level))) * hashSize)
	if pos < next { // 分裂过了的
		pos = posCRC16 % (int64(math.Pow(2.0, float64(level+1))) * hashSize)
	}
	return pos
}
