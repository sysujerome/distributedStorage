package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "example.com/kvstore"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedStorageServer
}

var db map[string]string

var (
	port = flag.Int("port", 50051, "The server port")
)

func (s *server) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
	// log.Printf("Received get key: %v", in.GetKey())
	value, found := db[in.GetKey()]
	if !found {
		return &pb.GetReply{Result: "not fount request: " + in.GetKey()}, nil
	}
	return &pb.GetReply{Result: "succeed! request: " + in.GetKey() + " get: " + value}, nil
}
func (s *server) Set(ctx context.Context, in *pb.SetRequest) (*pb.SetReply, error) {
	// log.Printf("Received get key: %v", in.GetKey())
	db[in.GetKey()] = in.GetValue()
	return &pb.SetReply{Result: "succeed insert " + in.GetKey() + " " + in.GetValue()}, nil
}
func (s *server) Del(ctx context.Context, in *pb.DelRequest) (*pb.DelReply, error) {
	// log.Printf("Received get key: %v", in.GetKey())
	delete(db, in.GetKey())
	_, found := db[in.GetKey()]
	if found {
		log.Fatal("del %v failed: ", in.GetKey())
	}
	return &pb.DelReply{Result: "succeed delete " + in.GetKey()}, nil
}

func main() {
	flag.Parse()
	db = make(map[string]string)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
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
