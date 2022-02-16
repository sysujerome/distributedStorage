package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "example.com/kvstore"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedStorageServer
}

var ctx = context.Background()

var rdb *redis.Client

var (
	port = flag.Int("port", 50052, "The server port")
)

func (s *server) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
	// log.Printf("Received get key: %v", in.GetKey())
	res, err := rdb.Get(ctx, in.GetKey()).Result()
	if err != nil {
		return &pb.GetReply{Result: "not fount request: " + in.GetKey()}, nil
	}
	return &pb.GetReply{Result: "succeed! request: " + in.GetKey() + " get: " + res}, nil
}
func (s *server) Set(ctx context.Context, in *pb.SetRequest) (*pb.SetReply, error) {
	// log.Printf("Received get key: %v", in.GetKey())
	res, err := rdb.Set(ctx, in.GetKey(), in.GetValue(), 0).Result()
	if err != nil {
		log.Fatal("get redis db failed: ", err)
	}
	return &pb.SetReply{Result: "succeed " + string(res)}, nil
}
func (s *server) Del(ctx context.Context, in *pb.DelRequest) (*pb.DelReply, error) {
	// log.Printf("Received get key: %v", in.GetKey())
	res, err := rdb.Del(ctx, in.GetKey()).Result()
	if err != nil {
		log.Fatal("get redis db failed: ", err)
	}
	return &pb.DelReply{Result: "succeed " + string(res)}, nil
}

func main() {
	flag.Parse()
	rdb = redis.NewClient(&redis.Options{
		Addr: "192.168.1.128:6379",
	})
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
