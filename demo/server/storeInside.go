package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	pb "example.com/kvstore"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedStorageServer
}

var db map[string]string
var next int

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
	return &pb.SetReply{Status: "OK"}, nil
}

func (s *server) Del(ctx context.Context, in *pb.DelRequest) (*pb.DelReply, error) {
	// log.Printf("Received get key: %v", in.GetKey())
	delete(db, in.GetKey())
	_, found := db[in.GetKey()]
	if found {
		return &pb.DelReply{Err: "failed"}, nil
	}
	return &pb.DelReply{Result: 1, Status: "OK"}, nil
}

func main() {
	flag.Parse()
	port := 0
	for i := 0; i < 5; i++ {
		port = 50050 + i
		go serve(port)
	}
	for {
		fmt.Println(time.Now().Format("15:04:05"))
		time.Sleep(time.Duration(1) * time.Minute)
	}
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

func split() {

}
