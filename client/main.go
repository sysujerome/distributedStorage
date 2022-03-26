package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	pb "storage/kvstore"
	"storage/util/common"
	commom "storage/util/common"
	crc16 "storage/util/crc16"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	serversAddress map[int64]string
	serversStatus  map[int64]string
	serversMaxKey  map[int64]int64
	next           int64
	level          int64
	hashSize       int64
	ip             = flag.String("ip", "", "the base server to get configure")
	port           = flag.Int64("port", 0, "the base server to get configure")
	conf           = flag.String("conf", "", "the defaulted configure file")
	statusCode     commom.Status
	errorCode      commom.Error
	serverStatus   common.ServerStatus
	version        int64
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	flag.Parse()
	initConf()

	var conns []*grpc.ClientConn
	for _, addr := range serversAddress {
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		check(err)
		conns = append(conns, conn)
	}

	// currentPath, err := os.Getwd()
	// check(err)
	// filename := filepath.Join(currentPath, "data/benchmark_data")
	// f, err := os.Open(filename)
	// check(err)
	// defer f.Close()

	// sc := bufio.NewScanner(f)
	// counter := 0
	// start := time.Now()
	// for sc.Scan() {
	// 	operation := strings.Split(sc.Text(), " ")
	// 	key := operation[1]
	// 	idx := hashFunc(key) % 1
	// 	getServe(operation, conns[idx])
	// 	counter++
	// }
	// duration := time.Since(start)
	// fmt.Printf("dealing with %d operations took %v Seconds\n", counter, duration.Seconds())
	// var operation, key, value string
	for {

		// fmt.Scanln(&operation, &key, &value)
		reader := bufio.NewReader(os.Stdin)
		data, _, _ := reader.ReadLine()
		operations := strings.Split(string(data), " ")

		switch operations[0] {
		case "scan":
			scan(operations)
		case "get":
			get(operations)
		case "set":
			set(operations)

		case "test":
			test()
		case "quit":
			for _, conn := range conns {
				conn.Close()
			}
			return
		}

	}

}

func scan(operations []string) {
	var conns []*grpc.ClientConn
	for idx := 0; idx < len(serversAddress); idx++ {
		addr := serversAddress[int64(idx)]
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		check(err)
		conns = append(conns, conn)
	}
	total := 0
	for idx := 0; idx < len(serversAddress); idx++ {
		server := serversAddress[int64(idx)]
		if len(operations) > 1 {
			if operations[1] == "conf" {
				c := pb.NewStorageClient(conns[0])
				ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Second))
				defer cancel()
				reply, err := c.GetConf(ctx, &pb.GetConfRequest{})
				check(err)
				fmt.Printf("conf: \n%s\n", reply.GetResult())
				servers := reply.GetServer()
				for _, server := range servers {
					fmt.Printf("idx : %d\n", server.Idx)
					fmt.Printf("address : %s\n", server.Address)
					fmt.Printf("maxKey: %d\n", server.MaxKey)
					fmt.Printf("status: %s\n\n", server.Status)
				}
				break
			}

			index, err := strconv.Atoi(operations[1])
			check(err)
			c := pb.NewStorageClient(conns[index])
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Second))
			defer cancel()
			reply, err := c.Scan(ctx, &pb.ScanRequest{})
			check(err)
			fmt.Printf("status: %s\n", serversStatus[int64(idx)])
			fmt.Printf("count: %d\n", reply.GetCount())
			fmt.Printf("result:\n%s\n\n", reply.GetResult())
			break
		}
		fmt.Printf("idx: %d\n", idx)
		fmt.Println(server)
		fmt.Printf("status: %s\n", serversStatus[int64(idx)])
		c := pb.NewStorageClient(conns[idx])
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Second))
		defer cancel()
		reply, err := c.Scan(ctx, &pb.ScanRequest{})
		check(err)
		count := reply.GetCount()
		fmt.Printf("count: %d\n\n", count)
		total += int(count)
	}
	fmt.Printf("Total : %d\n", total)
	for _, conn := range conns {
		conn.Close()
	}
}
func set(operations []string) {
	if len(operations) < 3 {
		// panic("too less agrs")
		fmt.Println("too less arguments")
		return
	}
	key := operations[1]
	value := operations[2]
	var conns []*grpc.ClientConn
	for idx := 0; idx < len(serversAddress); idx++ {
		addr := serversAddress[int64(idx)]
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		check(err)
		conns = append(conns, conn)
	}
	idx := hashFunc(key)
	fmt.Printf("idx: %d\n", idx)
	c := pb.NewStorageClient(conns[idx])
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Second))
	defer cancel()
	reply, err := c.Set(ctx, &pb.SetRequest{Key: key, Value: value})
	check(err)
	status := reply.GetStatus()
	result := reply.GetResult()
	fmt.Printf("status: %s\nresult: %s\n", status, result)
	if status != statusCode.Ok {
		fmt.Printf("error: %s\n", reply.GetErr())
	}
	for _, conn := range conns {
		conn.Close()
	}
}
func get(operations []string) {
	if len(operations) < 2 {
		// panic("too less agrs")
		fmt.Println("too less arguments")
		return
	}
	key := operations[1]
	var conns []*grpc.ClientConn
	for idx := 0; idx < len(serversAddress); idx++ {
		addr := serversAddress[int64(idx)]
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		check(err)
		conns = append(conns, conn)
	}

	idx := hashFunc(key)
	fmt.Printf("idx: %d\n", idx)
	c := pb.NewStorageClient(conns[idx])
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Second))
	defer cancel()
	reply, err := c.Get(ctx, &pb.GetRequest{Key: key})
	check(err)
	status := reply.GetStatus()
	result := reply.GetResult()
	fmt.Printf("status: %s\nresult: %s\n", status, result)
	if status != statusCode.Ok {
		fmt.Printf("error: %s\n", reply.GetErr())
	}
	for _, conn := range conns {
		conn.Close()
	}
}

func test() {

	var statusCode commom.Status
	var errorCode commom.Error
	statusCode.Init()
	errorCode.Init()

	serverAddress := make(map[int64]string)
	serverStatus := make(map[int64]string)
	serverMaxKey := make(map[int64]int64)

	address := serversAddress[0]
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	check(err)
	defer conn.Close()
	c := pb.NewStorageClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Minute))
	defer cancel()
	reply, err := c.GetConf(ctx, &pb.GetConfRequest{})
	check(err)
	if reply.GetStatus() != statusCode.Ok {
		fmt.Println("同步服务器配置出错...")
	}

	next = reply.GetNext()
	level = reply.GetLevel()
	hashSize = reply.GetHashSize()
	for _, server := range reply.GetServer() {
		idx := server.Idx
		serverAddress[idx] = server.Address
		serverMaxKey[idx] = server.MaxKey
		serverStatus[idx] = server.Status
	}
	version = reply.GetVersion()

	var conns []*grpc.ClientConn
	var clients []pb.StorageClient
	for idx := 0; idx < len(serversAddress); idx++ {
		addr := serversAddress[int64(idx)]
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		check(err)
		conns = append(conns, conn)
		c := pb.NewStorageClient(conn)
		clients = append(clients, c)
	}

	keys := []string{}
	values := []string{}

	path, err := os.Getwd()
	check(err)
	filename := filepath.Join(path, "./data/data2")
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		kv := strings.Split(sc.Text(), " ")
		keys = append(keys, kv[0])
		values = append(values, kv[1])
		// fmt.Println(kv[0], kv[1])
	}
	fmt.Println("set....")
	// SET
	start := time.Now()
	preSecond := time.Now()
	preIndex := 0
	// for i := 0; i < len(keys); i++ {
	// 	key := keys[i]
	// 	value := values[i]
	// 	// addr := serverAddress[idx]
	// 	// conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// 	// check(err)d
	// 	// defer conn.Close()
	// 	// c := pb.NewStorageClient(conn)

	// 	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Second))
	// 	defer cancel()
	// 	// 每次都要同步配置文件
	// 	// syncConf(clients[0], ctx)
	// 	idx := hashFunc(key)
	// 	reply, err := clients[idx].Set(ctx, &pb.SetRequest{Key: key, Value: value})
	// 	check(err)
	// 	if reply.GetVersion() != version {
	// 		syncConf(clients[idx], ctx)
	// 	}
	// 	if err != nil || reply.GetStatus() == statusCode.Failed {
	// 		fmt.Printf("%s status: %s\n", serversAddress[idx], reply.GetStatus())
	// 		panic(reply.GetErr())
	// 	}

	// 	// 如果出现moved状态码，说明集群状态已经改变，需要同步配置文件
	// 	if reply.GetStatus() == statusCode.Moved {
	// 		target := reply.GetTarget()
	// 		syncConf(clients[target], ctx)
	// 		reply1, err := clients[target].Set(ctx, &pb.SetRequest{Key: key, Value: value})
	// 		check(err)
	// 		if err != nil || reply1.GetStatus() == statusCode.Failed {
	// 			fmt.Printf("%s status: %s\n", serversAddress[target], reply1.GetStatus())
	// 			panic(reply1.GetErr())
	// 		}
	// 	}
	// 	if reply.GetVersion() != version {
	// 		syncConf(clients[idx], ctx)
	// 	}
	// 	// fmt.Println(reply.GetStatus())

	// 	// 按间隔输出qps
	// 	nextSecond := preSecond.Add(1 * time.Second)
	// 	now := time.Now()
	// 	if now.After(nextSecond) {
	// 		fmt.Printf("%d echo... per second\n", i-preIndex)
	// 		preIndex = i
	// 		preSecond = now
	// 	}
	// }
	elapse := time.Since(start)
	fmt.Printf("Set %d keys took %s\n", len(keys), elapse)

	fmt.Println("get....")
	// GET
	start = time.Now()
	preSecond = time.Now()
	preIndex = 0
	successNumber := 0
	wrongNumber := 0
	syncConf(clients[0], ctx)
	// start = time.Now()
	// for i := 0; i < len(keys); i++ {
	// 	// if i%10 == 0 {
	// 	// 	fmt.Printf("%d epoch...\n", i)
	// 	// }
	// 	key := keys[i]
	// 	// value := values[i]
	// 	idx := hashFunc(key)
	// 	// found := false

	// 	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Second))
	// 	defer cancel()
	// 	reply, _ := clients[idx].Get(ctx, &pb.GetRequest{Key: key})
	// 	// if reply.GetStatus() == statusCode.Ok && reply.GetResult() == values[i] {
	// 	if reply.GetStatus() == statusCode.Ok {
	// 		successNumber++
	// 	}
	// 	// 按间隔输出qps
	// 	nextSecond := preSecond.Add(1 * time.Second)
	// 	now := time.Now()
	// 	if now.After(nextSecond) {
	// 		fmt.Printf("%d echo... per second\n", i-preIndex)
	// 		preIndex = i
	// 		preSecond = now
	// 	}
	// }
	elapse = time.Since(start)
	fmt.Printf("Get %d keys took %s\n", len(keys), elapse)

	fmt.Printf("Total number : %d\n", len(keys))
	fmt.Printf("Success number : %d\n", successNumber)
	fmt.Printf("Wrong number : %d\n", wrongNumber)

	for idx := 0; idx < len(serversAddress); idx++ {
		conns[idx].Close()
	}
	return

	// DEL
	start = time.Now()
	preSecond = time.Now()
	// fmt.Println("")
	fmt.Println("del....")
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
		if err != nil {
			panic(err)
		}
		if reply.GetStatus() != statusCode.Ok {
			panic(reply.GetErr())
		}
		if reply.GetResult() != 1 {
			log.Fatalf("Delete key %s failed!", key)
		}
		// 按间隔输出qps
		nextSecond := preSecond.Add(1 * time.Second)
		now := time.Now()
		if now.After(nextSecond) {
			fmt.Printf("%d echo... per second\n", i-preIndex)
			preIndex = i
			preSecond = now
		}
	}
	elapse = time.Since(start)
	fmt.Printf("Del %d keys took %s\n", len(keys), elapse)

	// GET
	fmt.Println("get....")
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
func getServe(operation []string, conn *grpc.ClientConn) {

	c := pb.NewStorageClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Second))
	defer cancel()

	opt := operation[0]
	switch opt {
	case "get":
		reply, err := c.Get(ctx, &pb.GetRequest{Key: operation[1]})
		check(err)
		fmt.Printf("%s", reply.GetStatus())
	case "set":
		for len(operation) < 3 {
			operation = append(operation, "dafaskjdhf")
		}
		reply, err := c.Set(ctx, &pb.SetRequest{Key: operation[1], Value: operation[2]})
		check(err)
		fmt.Printf("%s", reply.GetStatus())
	case "del":
		reply, err := c.Del(ctx, &pb.DelRequest{Key: operation[1]})
		check(err)
		fmt.Printf("%s", reply.GetStatus())
	}

}

func initConf() {

	serversAddress = make(map[int64]string)
	serversStatus = make(map[int64]string)
	serversMaxKey = make(map[int64]int64)
	address := ""
	errorCode.Init()
	statusCode.Init()
	serverStatus.Init()
	defer flag.Parse()
	curPath, err := os.Getwd()
	check(err)
	if *conf != "" {
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

			ip := ShardNodeConfs["ip"].(string)
			port := int64(ShardNodeConfs["base_port"].(float64))
			address := fmt.Sprintf("%s:%d", ip, port)
			serversAddress[idx] = address
			serversStatus[idx] = ShardNodeConfs["status"].(string)
			serversMaxKey[idx] = int64(ShardNodeConfs["max_key"].(float64))
			// fmt.Println(idx, ip, port)
		}
		address = serversAddress[0]
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
	syncConf(c, ctx)
	// printConf()
}

func syncConf(c pb.StorageClient, ctx context.Context) {

	reply, err := c.GetConf(ctx, &pb.GetConfRequest{})
	check(err)
	if reply.GetStatus() != statusCode.Ok {
		fmt.Println(reply.GetStatus())
		fmt.Println(statusCode.Ok)
		fmt.Println("同步服务器配置出错...")
	}
	next = reply.GetNext()
	level = reply.GetLevel()
	hashSize = reply.GetHashSize()
	for _, server := range reply.GetServer() {
		idx := server.Idx
		serversAddress[int64(idx)] = server.Address
		serversMaxKey[int64(idx)] = server.MaxKey
		serversStatus[int64(idx)] = server.Status
	}
	version = reply.GetVersion()
}

func printConf() {
	for idx, addr := range serversAddress {
		fmt.Printf("%d %s\t", idx, addr)
	}
	fmt.Println()

	for idx, maxKey := range serversMaxKey {
		fmt.Printf("%d %d\t", idx, maxKey)
	}
	fmt.Println()

	for idx, status := range serversStatus {
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
	// return 1
}
