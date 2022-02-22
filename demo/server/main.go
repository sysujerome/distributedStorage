package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	pb "example.com/kvstore"
	crc16 "example.com/util/crc16"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type server struct {
	pb.UnimplementedStorageServer
}

var (
	db       map[string]string
	next     int
	port     int
	maxCount int
	curCount int
	shard_node_name 
)

func init() {
	maxCount = 10000
	curCount = 0
}

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
	addr := "192.168.1.128" + fmt.Sprintf(":%d", in.GetPort())
	log.Println(addr)
	count := split(addr)
	return &pb.SplitReply{Status: "OK", Result: count}, nil
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
	flag.Parse()
	port, err := strconv.Atoi(os.Args[1])
	check(err)
	serve(port)
}

func serve(port int) {
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
		if crc16.Checksum([]byte(key), crc16.IBMTable)%2 == 1 { // 待修改
			delete(db, key)
			reply1, _ := c.Set(ctx, &pb.SetRequest{Key: key, Value: value})
			if reply1.GetStatus() != "OK" {
				fmt.Printf("\"%v\"", reply1.GetStatus())
			}
			count++
		}
		// fmt.Printf("%v:%v", key, value)
		// reply, err := c.Set(ctx, &pb.SetRequest{Key: operation[1], Value: operation[2]})
	}
	return int32(count)
}

func initConf() {

}
