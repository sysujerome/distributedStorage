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
	"strconv"
	"time"

	pb "example.com/kvstore"
	// "example.com/util/crc16"
	crc16 "example.com/util/crc16"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type server struct {
	pb.UnimplementedStorageServer
}

var (
	db            map[string]string
	next          = flag.Int64("next", 0, "the spilt node point64er")
	maxCount      = flag.Int64("maxCount", 10000, "the maximum count of the key when to split")
	curCount      = flag.Int64("curCount", 0, "the count of the key at the beginning")
	shardNodeName = flag.String("shard_node_name", "shard_node_0", "incade the node name")
	shardIdx      = flag.Int64("shard_idx", 0, "the index of the shard node")
	ip            = flag.String("ip", "192.168.1.128", "indicate the ip of the node")
	basePort      = flag.Int64("base_port", 50050, "indicate the port of the node")
	conf          = flag.String("conf", "conf.json", "the path of configure file")
	usingSize     = flag.Int64("using_size", 4, "the free machine to add the data of the splitting machine")
	maxSize       = flag.Int64("max_size", 8, "the maximum of the count of machine")
	hashSize      = flag.Int64("hash_size", 4, "the initial count of hash number")
	level         = flag.Int64("level", 4, "the initial level of hash function")
	servers       map[int64]string
	maxKey        map[int64]int64
	status        map[int64]int64
)

// var servers []string

func (s *server) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
	// log.Printf("Received get key: %v", in.GetKey())
	value, found := db[in.GetKey()]
	if !found {
		return &pb.GetReply{Result: "", Status: "Not found!"}, nil
	}
	return &pb.GetReply{Result: value, Status: "OK"}, nil
}

func (s *server) Set(ctx context.Context, in *pb.SetRequest) (*pb.SetReply, error) {
	// log.Printf("Received get key: %v", in.GetKey())
	db[in.GetKey()] = in.GetValue()
	if db[in.GetKey()] != in.GetValue() {
		return &pb.SetReply{Status: "failed!!!"}, nil
	}
	return &pb.SetReply{Status: "OK"}, nil
}

func (s *server) Del(ctx context.Context, in *pb.DelRequest) (*pb.DelReply, error) {
	// log.Printf("Received get key: %v", in.GetKey())
	delete(db, in.GetKey())
	_, found := db[in.GetKey()]
	if found {
		return &pb.DelReply{Err: "failed", Status: "Not found!"}, nil
	}
	return &pb.DelReply{Result: 1, Status: "OK"}, nil
}

func (s *server) Split(ctx context.Context, in *pb.SplitRequest) (*pb.SplitReply, error) {
	// log.Printf("Received get key: %v", in.GetKey())
	if *usingSize >= *maxSize {
		*next = *maxSize
		return &pb.SplitReply{Status: "failed", Result: 0, Err: "the count of machine using haved fulled!"}, nil
	}
	addr := "192.168.1.128" + fmt.Sprintf(":%d", *basePort+*usingSize)
	*usingSize++
	log.Println(addr)
	count := split(addr)
	*next++
	if *next == int64(math.Pow(2, float64(*level)))**hashSize {
		*next = 0
		*level++
	}
	return &pb.SplitReply{Status: "OK", Result: count}, nil
}

func (s *server) SyncConf(ctx context.Context, in *pb.SyncRequest) (*pb.SyncReply, error) {
	*next = in.GetNext()
	*level = in.GetLevel()
	syncConf(in.GetBeginLevel())
	return &pb.SyncReply{Status: "OK"}, nil
}

// func (c *storageClient) Scan(ctx context.Context, in *ScanRequest, opts ...grpc.CallOption) (*ScanReply, error)
func (s *server) Scan(ctx context.Context, in *pb.ScanRequest) (*pb.ScanReply, error) {
	// log.Printf("Received get key: %v", in.GetKey())
	// addr := "192.168.1.128:" + fmt.Sprintf(":%d", in.GetPort())
	// log.Println(addr)
	result := ""
	count := 0
	for key, value := range db {
		result += key + ":" + value + "\n"
		count++
	}
	result = strconv.Itoa(count) + "\n" + result
	result = "OK" + result
	return &pb.ScanReply{Status: "OK", Result: result}, nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	initConf()
	serve(*basePort)
}

func printConf() {
	fmt.Printf("db              %s\n", db)
	fmt.Printf("servers         %s\n", servers)
	fmt.Printf("next            %d\n", *next)
	fmt.Printf("maxCount        %d\n", *maxCount)
	fmt.Printf("curCount        %d\n", *curCount)
	fmt.Printf("shard_node_name %s\n", *shardNodeName)
	fmt.Printf("shard_idx       %d\n", *shardIdx)
	fmt.Printf("ip              %s\n", *ip)
	fmt.Printf("basePort        %d\n", *basePort)
	fmt.Printf("conf            %s\n", *conf)
}

func serve(port int64) {
	db = make(map[string]string)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
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

func split(target string) int32 {
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	check(err)
	c := pb.NewStorageClient(conn)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Minute))
	defer cancel()

	count := 0
	for key, value := range db {
		if hashFunc(key) != *shardIdx { // 移到新的节点
			delete(db, key)
			reply1, _ := c.Set(ctx, &pb.SetRequest{Key: key, Value: value})
			if reply1.GetStatus() != "OK" {
				// fmt.Printf("\"%v\"", reply1.GetStatus())
				panic(reply1.GetStatus())
			}
			count++
		}
	}
	return int32(count)
}

func initConf() {
	flag.Parse()
	db = make(map[string]string)
	servers = make(map[int64]string)
	fmt.Printf("%s\n", *shardNodeName)
	defer flag.Parse()
	curPath, err := os.Getwd()
	check(err)
	filePath := filepath.Join(curPath, *conf)
	//检查文件是否存在
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		panic("配置文件路径出错")
	}
	fileData, err := ioutil.ReadFile(filePath)
	var configure map[string]interface{}
	json.Unmarshal(fileData, &configure)
	confs := configure["shard_confs"].([]interface{}) // shard_confs
	for _, v := range confs {
		value := v.(map[string]interface{})
		index := int64(value["shard_idx"].(float64))
		shard_node_confs := value["shard_node_confs"].(map[string]interface{})
		shard_content, found := shard_node_confs[*shardNodeName].(map[string]interface{})
		if found {
			*shardIdx = index
			*basePort = int64(shard_content["base_port"].(float64))
			*ip = shard_content["ip"].(string)
			*maxCount = int64(shard_content["max_key"].(float64))
		}
		for _, conf := range shard_node_confs {
			detail := conf.(map[string]interface{})
			thisIP := detail["ip"].(string)
			thisPort := detail["base_port"].(float64)
			address := thisIP + ":" + strconv.Itoa(int(thisPort))
			servers[index] = address
		}
	}
	*level = 0
}

func syncConf(beginIdx int64) { // beginIdx指的是发起同步的node
	getNextIdx := func(index int64, max int) int64 {
		nextIndex := index + 1
		if nextIndex == int64(max) {
			nextIndex = 0
		}
		return nextIndex
	}
	nextIndex := getNextIdx(*shardIdx, len(servers))
	for {
		if nextIndex == *shardIdx || nextIndex == beginIdx { //同步到发起节点或者同步到自身则停止
			return
		}
		address := servers[nextIndex]
		_, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			nextIndex = getNextIdx(nextIndex, len(servers))
			continue
		}
		conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		check(err)
		c := pb.NewStorageClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Hour))
		defer cancel()
		reply, err := c.SyncConf(ctx, &pb.SyncRequest{Next: *next, Level: *level})
		check(err)
		if reply.GetStatus() != "OK" {
			nextIndex = getNextIdx(nextIndex, len(servers))
			continue
		}
		break
	}

}

func hashFunc(key string) int64 {
	posCRC16 := int64(crc16.Checksum([]byte(key), crc16.IBMTable))
	pos := posCRC16 % (int64(math.Pow(2.0, float64(*level))) * *hashSize)
	if pos < *next { // 分裂过了的
		pos = posCRC16 % (int64(math.Pow(2.0, float64(*level+1))) * *hashSize)
	}
	return pos
}
