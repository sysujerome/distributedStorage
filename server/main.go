package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net"
	"os"
	"path/filepath"
	"reflect"
	pb "storage/kvstore"
	commom "storage/util/common"
	crc16 "storage/util/crc16"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	db            map[string]string
	serverAddress map[int64]string
	serverStatus  map[int64]string
	serverMaxKey  map[int64]int64
	shardIdx      = flag.Int64("shard_idx", 0, "the index of this shard node")
	next          = flag.Int64("next", 0, "the spilt node pointer")
	level         = flag.Int64("level", 0, "the initial level of hash function")
	hashSize      = flag.Int64("hash_size", 4, "the initial count of hash number")
	conf          = flag.String("conf", "conf.json", "the defaulted configure file")
	statusCode    commom.Status
	errorCode     commom.Error
)

type server struct {
	pb.UnimplementedStorageServer
}

func (s *server) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
	// log.Printf("Received get key: %v", in.GetKey())
	value, found := db[in.GetKey()]
	if !found {
		return &pb.GetReply{Result: "", Err: errorCode.NotFound, Status: statusCode.Failed}, nil
	}
	return &pb.GetReply{Result: value, Status: statusCode.Ok}, nil
}

func (s *server) Set(ctx context.Context, in *pb.SetRequest) (*pb.SetReply, error) {
	// log.Printf("Received get key: %v", in.GetKey())
	db[in.GetKey()] = in.GetValue()
	if db[in.GetKey()] != in.GetValue() {
		return &pb.SetReply{Status: "failed!!!"}, nil
	}
	return &pb.SetReply{Status: statusCode.Ok}, nil
}

func (s *server) Del(ctx context.Context, in *pb.DelRequest) (*pb.DelReply, error) {
	// log.Printf("Received get key: %v", in.GetKey())
	delete(db, in.GetKey())
	_, found := db[in.GetKey()]
	if found {
		return &pb.DelReply{Err: errorCode.NotFound, Status: statusCode.Failed}, nil
	}
	return &pb.DelReply{Result: 1, Status: statusCode.Ok}, nil
}

func (s *server) Split(ctx context.Context, in *pb.SplitRequest) (*pb.SplitReply, error) {
	// log.Printf("Received get key: %v", in.GetKey())
	secondIdx := int64(math.Pow(2, float64(*level))*float64(*hashSize) + float64(*shardIdx))
	if secondIdx >= int64(len(serverAddress)) {
		return &pb.SplitReply{Status: statusCode.Failed, Err: "The server has fulled!"}, nil
	}

	split := func(target string) int64 { // 内部元素分裂
		conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
		check(err)
		c := pb.NewStorageClient(conn)
		// Contact the server and print out its response.
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Second))
		defer cancel()

		count := 0
		for key, value := range db {
			if crc16.Checksum([]byte(key), crc16.IBMTable)%2 == 1 { // 待修改
				delete(db, key)
			}
			reply1, _ := c.Set(ctx, &pb.SetRequest{Key: key, Value: value})
			if reply1.GetStatus() != "OK" {
				panic(reply1.GetStatus())
			}
			count++
		}
		return int64(count)
	}

	count := split(serverAddress[secondIdx])
	*next++
	if *next == int64(math.Pow(2, float64(*level)))**hashSize {
		*next = 0
		*level++
	}
	return &pb.SplitReply{Status: "OK", Result: count}, nil
}

func (s *server) Scan(ctx context.Context, in *pb.ScanRequest) (*pb.ScanReply, error) {
	count := len(db)
	result := serverAddress[*shardIdx] + "\n"
	result += strconv.Itoa(count) + "\n"
	result += fmt.Sprintf("%s", db)
	return &pb.ScanReply{Result: result, Status: statusCode.Ok}, nil
}

func (s *server) SyncConf(ctx context.Context, in *pb.SyncConfRequest) (*pb.SyncConfReply, error) {
	if in.GetBegin() == *shardIdx {
		return &pb.SyncConfReply{Status: statusCode.Ok, Result: "The sync idx has come to begin!"}, nil
	}

	*next = in.GetNext()
	*level = in.GetLevel()
	*hashSize = in.GetHashSize()
	for k := range serverAddress {
		delete(serverAddress, k)
		delete(serverAddress, k)
		delete(serverAddress, k)
	}
	for _, serv := range in.GetServer() {
		fmt.Printf("%s\n", reflect.TypeOf(serv))
		serverAddress[serv.Idx] = serv.Address
		serverMaxKey[serv.Idx] = serv.MaxKey
		serverStatus[serv.Idx] = serv.Status
	}
	success := false
	target := *shardIdx + 1
	for !success {
		serv := serverAddress[target]
		conn, err := grpc.Dial(serv, grpc.WithTransportCredentials(insecure.NewCredentials()))
		check(err)
		c := pb.NewStorageClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Second))
		defer cancel()
		reply, err := c.SyncConf(ctx, in)
		check(err)
		if reply.GetStatus() == statusCode.Ok {
			break
		}
		target++
		if target == in.GetBegin() {
			success = true
		}
		if target == int64(len(serverAddress)) {
			target = 0
		}
	}
	return &pb.SyncConfReply{Status: statusCode.Ok}, nil
}

func (s *server) GetConf(ctx context.Context, in *pb.GetConfRequest) (*pb.GetConfReply, error) {
	var reply pb.GetConfReply
	reply.HashSize = *hashSize
	reply.Level = *level
	reply.Next = *next
	reply.Status = statusCode.Ok
	for idx, addr := range serverAddress {
		var server pb.GetConfReply_ServConf
		server.Address = addr
		server.Idx = idx
		server.MaxKey = serverMaxKey[idx]
		server.Status = serverStatus[idx]
		reply.Server = append(reply.Server, &server)
	}
	return &reply, nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func initConf() {
	flag.Parse()

	db = make(map[string]string)
	serverAddress = make(map[int64]string)
	serverStatus = make(map[int64]string)
	serverMaxKey = make(map[int64]int64)
	defer flag.Parse()
	curPath, err := os.Getwd()
	check(err)
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
	*level = 0
}

func main() {
	flag.Parse()
	initConf()
	// printConf()
	serve()
}

func printConf() {
	fmt.Printf("db:\n%s\n", db)
	// fmt.Printf("servers:\n%s\n", serverAddress)
	for idx, addr := range serverAddress {
		fmt.Printf("%d %s\t", idx, addr)
	}
	fmt.Println()
	// fmt.Printf("status:\n%s\n", serverStatus)
	for idx, maxKey := range serverMaxKey {
		fmt.Printf("%d %d\t", idx, maxKey)
	}
	fmt.Println()

	// fmt.Printf("maxKey:\n%s\n", serverMaxKey)
	for idx, status := range serverStatus {
		fmt.Printf("%d %s\t", idx, status)
	}
	fmt.Println()
	fmt.Printf("shardIdx:\n%d\n", *shardIdx)
	fmt.Printf("next:\n%d\n", *next)
	fmt.Printf("level:\n%d\n", *level)
	fmt.Printf("hashSize:\n%d\n", *hashSize)
	fmt.Printf("conf:\n%s\n", *conf)
}

func serve() {
	addr := serverAddress[*shardIdx]
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterStorageServer(s, &server{}) //StorageServer
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
