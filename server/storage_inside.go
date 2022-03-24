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
	pb "storage/kvstore"
	common "storage/util/common"
	crc16 "storage/util/crc16"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	db             map[string]string
	operations     [][]string
	serversAddress map[int64]string
	serversStatus  map[int64]string
	serversMaxKey  map[int64]int64
	shardIdx       = flag.Int64("shard_idx", 0, "the index of this shard node")
	next           = flag.Int64("next", 0, "the spilt node pointer")
	level          = flag.Int64("level", 0, "the initial level of hash function")
	hashSize       = flag.Int64("hash_size", 4, "the initial count of hash number")
	conf           = flag.String("conf", "conf.json", "the defaulted configure file")
	statusCode     common.Status
	errorCode      common.Error
	serverStatus   common.ServerStatus
	canSplit       bool
	// mutex          sync.Mutex
	overflowDB *leveldb.DB
	spliting int64
)

type server struct {
	pb.UnimplementedStorageServer
}

func (s *server) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
	if serversStatus[*shardIdx] == serverStatus.Sleep {
		return &pb.GetReply{Result: "", Status: statusCode.Failed, Err: errorCode.NotWorking}, nil
	}

	if atomic.LoadInt64(&spliting) > 0 {
		// operations = append(operations, []string{"get", in.GetKey()})
		value, found := db[in.GetKey()]
		if found {
			return &pb.GetReply{Result: value, Status: statusCode.Ok}, nil
		} else {
			// 从溢出表中查找
			// data, err := overflowDB.Get([]byte("key"), nil)
			// if (err != nil) {
			// 	return &pb.GetReply{Result: string(data), Status: statusCode.Ok}, nil
			// } else {
				secondIdx := int64(math.Pow(2, float64(*level)))**hashSize + *shardIdx
				if int(secondIdx) >= len(serversAddress) {
					return &pb.GetReply{Result: "", Status: statusCode.Failed, Err: errorCode.NotFound}, nil
				}
				target := serversAddress[secondIdx]
				conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
				check(err)
				defer conn.Close()
				c := pb.NewStorageClient(conn)
				// Contact the server and print out its response.
				ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Second))
				defer cancel()
	
				reply, err := c.Get(ctx, &pb.GetRequest{Key: in.GetKey()})
				check(err)
				if reply.GetStatus() == statusCode.Ok {
					return &pb.GetReply{Result: reply.GetResult(), Status: statusCode.Ok}, nil
				} else {
					return &pb.GetReply{Result: "", Status: statusCode.Failed, Err: errorCode.NotFound}, nil
				}
			// }
		}
		
		return &pb.GetReply{Result: "", Status: statusCode.Stored, Err: errorCode.Stored}, nil
	}
	value, found := db[in.GetKey()]
	if !found {
		// data, err := overflowDB.Get([]byte(in.GetKey()), nil)
		// if err != nil {
		// 	return &pb.GetReply{Result: "", Err: errorCode.NotFound, Status: statusCode.Failed}, nil
		// }
		return &pb.GetReply{Result: string(data), Err: "", Status: statusCode.Ok}, nil
	}
	return &pb.GetReply{Result: value, Status: statusCode.Ok}, nil
}

func (s *server) Set(ctx context.Context, in *pb.SetRequest) (*pb.SetReply, error) {
	if serversStatus[*shardIdx] == serverStatus.Sleep {
		return &pb.SetReply{Result: "", Status: statusCode.Failed, Err: errorCode.NotWorking}, nil
	}

	if atomic.LoadInt64(&spliting) > 0 {
		// time.Sleep(time.Millisecond)
		// fmt.Printf("-")
		operations = append(operations, []string{"set", in.GetKey(), in.GetValue()})
		return &pb.SetReply{Result: "", Status: statusCode.Stored, Err: errorCode.Stored}, nil
	}
	// }
	
	if len(db) == int(serversMaxKey[*shardIdx]) {
		// 继续放在本地db上
		db[in.GetKey()] = in.GetValue()
		if canSplit {
			split()
		}
	}
	db[in.GetKey()] = in.GetValue()
	if db[in.GetKey()] != in.GetValue() {
		return &pb.SetReply{Result: "", Status: statusCode.Failed, Err: errorCode.NotDefined}, nil
	}
	return &pb.SetReply{Status: statusCode.Ok}, nil
}

func (s *server) Del(ctx context.Context, in *pb.DelRequest) (*pb.DelReply, error) {
	if serversStatus[*shardIdx] == serverStatus.Sleep {
		return &pb.DelReply{Result: 0, Status: statusCode.Failed, Err: errorCode.NotWorking}, nil
	}

	// log.Printf("Received get key: %v", in.GetKey())
	// if serversStatus[*shardIdx] == serverStatus.Spliting {
	if atomic.LoadInt64(&spliting) > 0 {
		operations = append(operations, []string{"del", in.GetKey()})
		return &pb.DelReply{Result: 0, Status: statusCode.Stored, Err: errorCode.Stored}, nil
	}
	delete(db, in.GetKey())
	_, found := db[in.GetKey()]
	if found {
		return &pb.DelReply{Result: 0, Status: statusCode.Failed, Err: errorCode.NotDefined}, nil
	}
	return &pb.DelReply{Result: 1, Status: statusCode.Ok}, nil
}

func split() {
	if *next >= int64(len(serversAddress)) {
		fmt.Println("分裂失败, 节点满了。。。")
		return
	}
	addr := serversAddress[*next]
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	check(err)
	defer conn.Close()
	c := pb.NewStorageClient(conn)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Minute))
	defer cancel()

	reply, err := c.Split(ctx, &pb.SplitRequest{})
	check(err)
	if reply.GetFull() {
		fmt.Println("分裂失败, 节点满了。。。")
	}
}

func (s *server) Split(ctx context.Context, in *pb.SplitRequest) (*pb.SplitReply, error) {
	// log.Printf("Received get key: %v", in.GetKey())

	splitInside := func() { // 内部元素分裂
		fmt.Printf("%s : 节点分裂中....\n", serversAddress[*shardIdx])
		secondIdx := int64(math.Pow(2, float64(*level)))**hashSize + *shardIdx
		if secondIdx >= int64(len(serversAddress)) {
			fmt.Printf("%s : 节点数满了....\n", serversAddress[*shardIdx])
			canSplit = false
			return
		}

		*next++
		if *next == int64(math.Pow(2, float64(*level)))**hashSize {
			*next = 0
			*level++
		}	
		syncConf()
		target := serversAddress[secondIdx]
		conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
		check(err)
		defer conn.Close()
		c := pb.NewStorageClient(conn)
		// Contact the server and print out its response.
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Second))
		defer cancel()


		reply, err := c.WakeUp(ctx, &pb.WakeRequest{})
		check(err)
		if reply.GetStatus() != statusCode.Ok {
			// fmt.Printf("%s %s : 节点唤醒失败....\n", serversAddress[*shardIdx], serversAddress[secondIdx])
			return
		}
		syncConf()

		thisAddress := serversAddress[*shardIdx]
		thisConn, err := grpc.Dial(thisAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		check(err)
		defer thisConn.Close()
		thisClient := pb.NewStorageClient(thisConn)
		// Contact the server and print out its response.

		fmt.Printf("%s : 开始分裂....\n", serversAddress[secondIdx])
		count := 0
		for key, value := range db {
			if hashFunc(key) == secondIdx {
				delete(db, key)
				reply1, err := c.Set(ctx, &pb.SetRequest{Key: key, Value: value})
				check(err)
				if reply1.GetStatus() != statusCode.Ok {
					fmt.Println(reply1.GetStatus())
					fmt.Println("出错。。。")
					panic(reply1.GetErr())
				}
				count++
			}
		}
		// 遍历溢出表的内容
		iter := overflowDB.NewIterator(nil, nil)
		for iter.Next() {
			key := string(iter.Key())
			value := string(iter.Value())
			curIDX := hashFunc(string(key))
			err := overflowDB.Delete([]byte(key), nil)
			check(err)
			
			if curIDX == secondIdx {
				reply1, err := c.Set(ctx, &pb.SetRequest{Key: key, Value: value})
				check(err)
				if reply1.GetStatus() != statusCode.Ok {
					fmt.Println(reply1.GetStatus())
					fmt.Println("出错。。。")
					panic(reply1.GetErr())
				}
				count++
			} else {
				reply2, err := thisClient.Set(ctx, &pb.SetRequest{Key: key, Value: value})
				check(err)
				if reply2.GetStatus() != statusCode.Ok {
					fmt.Println(reply2.GetStatus())
					fmt.Println("出错。。。")
					panic(reply2.GetErr())
				}
			}
		}

		serversStatus[*shardIdx] = serverStatus.Working
		syncConf()
		fmt.Printf("%s : 分裂完成....\n", serversAddress[secondIdx])
		// 处理存起来的操作
		fmt.Println("处理存起来的操作...")
		for _, operation := range operations {
			switch operation[0] {
			case "set":
				idx := hashFunc(operation[1])
				if idx == *shardIdx {
					db[operation[1]] = db[operation[2]]
				} else {
					target := serversAddress[idx]
					conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
					check(err)
					defer conn.Close()
					c := pb.NewStorageClient(conn)
					reply, err := c.Set(ctx, &pb.SetRequest{Key: operation[1], Value: operation[2]})
					if err != nil || reply.GetStatus() != statusCode.Ok {
						panic(reply.GetStatus())
					}
				}
				// db[operation[1]] = operation[2]
			case "del":
				// delete(db, operation[1])
				idx := hashFunc(operation[1])
				if idx == *shardIdx {
					db[operation[1]] = db[operation[2]]
				} else {
					target := serversAddress[idx]
					conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
					check(err)
					defer conn.Close()
					c := pb.NewStorageClient(conn)
					c.Del(ctx, &pb.DelRequest{Key: operation[1]})
				}
			}
		}
		serversStatus[*shardIdx] = serverStatus.Working
		return 
	}

	// mutex.Lock()
	serversStatus[*shardIdx] = serverStatus.Spliting
	atomic.AddInt64(&spliting, 1)
	// mutex.Unlock()
	go splitInside()

	// mutex.Lock()
	serversStatus[*shardIdx] = serverStatus.Working
	atomic.AddInt64(&spliting, -1)
	// mutex.Unlock()
	// if full {
	// 	canSplit = false
	// 	syncConf()
	// 	return &pb.SplitReply{Status: statusCode.Ok, Result: 0, Full: true}, nil
	// }

	return &pb.SplitReply{Status: statusCode.Ok}, nil
}


func (s *server) Scan(ctx context.Context, in *pb.ScanRequest) (*pb.ScanReply, error) {
	count := len(db)
	result := serversAddress[*shardIdx] + "\n"
	result += fmt.Sprintf("%s", db)
	return &pb.ScanReply{Result: result, Status: statusCode.Ok, Count: int64(count)}, nil
}

func syncConf() {
	request := pb.SyncConfRequest{}
	request.Begin = *shardIdx
	request.HashSize = *hashSize
	request.Level = *level
	request.Next = *next
	request.Server = make([]*pb.SyncConfRequest_ServConf, 0)
	request.CanSplit = canSplit
	for idx, addr := range serversAddress {
		var server pb.SyncConfRequest_ServConf
		server.Address = addr
		server.Idx = idx
		server.MaxKey = serversMaxKey[idx]
		server.Status = serversStatus[idx]
		request.Server = append(request.Server, &server)
	}
	target := *shardIdx + 1
	if target == int64(len(serversAddress)) {
		target = 0
	}
	serv := serversAddress[target]
	conn, err := grpc.Dial(serv, grpc.WithTransportCredentials(insecure.NewCredentials()))
	check(err)
	c := pb.NewStorageClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Second))
	defer cancel()
	reply, err := c.SyncConf(ctx, &request)
	check(err)
	if reply.GetStatus() != statusCode.Ok {
		fmt.Printf("%s : 同步失败。。。\n", serversAddress[*shardIdx])
	}
}

func (s *server) SyncConf(ctx context.Context, in *pb.SyncConfRequest) (*pb.SyncConfReply, error) {
	if in.GetBegin() == *shardIdx {
		return &pb.SyncConfReply{Status: statusCode.Ok, Result: "The sync idx has come to begin!"}, nil
	}
	fmt.Printf("%s : 节点同步中...\n", serversAddress[*shardIdx])
	*next = in.GetNext()
	*level = in.GetLevel()
	*hashSize = in.GetHashSize()
	canSplit = in.GetCanSplit()
	for k := range serversAddress {
		delete(serversAddress, k)
		delete(serversAddress, k)
		delete(serversAddress, k)
	}
	for _, serv := range in.GetServer() {
		serversAddress[serv.Idx] = serv.Address
		serversMaxKey[serv.Idx] = serv.MaxKey
		serversStatus[serv.Idx] = serv.Status
	}
	success := false
	target := *shardIdx + 1
	if target == int64(len(serversAddress)) {
		target = 0
	}
	for !success {
		serv := serversAddress[target]
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
		if target == int64(len(serversAddress)) {
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
	for idx, addr := range serversAddress {
		var server pb.GetConfReply_ServConf
		server.Address = addr
		server.Idx = idx
		server.MaxKey = serversMaxKey[idx]
		server.Status = serversStatus[idx]
		reply.Server = append(reply.Server, &server)
	}
	return &reply, nil
}

func (s *server) WakeUp(ctx context.Context, in *pb.WakeRequest) (*pb.WakeReply, error) {
	if serversStatus[*shardIdx] != serverStatus.Sleep {
		return &pb.WakeReply{Status: statusCode.Failed, Result: "The server is not sleep..."}, nil
	}
	serversStatus[*shardIdx] = serverStatus.Working
	syncConf()
	return &pb.WakeReply{Status: statusCode.Ok}, nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func initConf() {
	flag.Parse()

	db = make(map[string]string)
	overflowDB, _ = leveldb.OpenFile("./overflow.db", nil)
	// check(err)
	operations = make([][]string, 0)
	serversAddress = make(map[int64]string)
	serversStatus = make(map[int64]string)
	serversMaxKey = make(map[int64]int64)
	canSplit = true
	atomic.StoreInt64(&spliting, 0)
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
		ip := ShardNodeConfs["ip"].(string)
		port := int64(ShardNodeConfs["base_port"].(float64))
		address := fmt.Sprintf("%s:%d", ip, port)
		serversAddress[idx] = address
		serversStatus[idx] = ShardNodeConfs["status"].(string)
		serversMaxKey[idx] = int64(ShardNodeConfs["max_key"].(float64))
		// fmt.Println(idx, ip, port)
	}
	*level = 0
	statusCode.Init()
	errorCode.Init()
	serverStatus.Init()
}

func hashFunc(key string) int64 {
	posCRC16 := int64(crc16.Checksum([]byte(key), crc16.IBMTable))
	pos := posCRC16 % (int64(math.Pow(2.0, float64(*level))) * *hashSize)
	if pos < *next { // 分裂过了的
		pos = posCRC16 % (int64(math.Pow(2.0, float64(*level+1))) * *hashSize)
	}
	return pos
	// return 1
}

func main() {
	flag.Parse()
	initConf()
	// printConf()
	serve()
}

func printConf() {
	fmt.Printf("db:\n%s\n", db)
	// fmt.Printf("servers:\n%s\n", serversAddress)
	for idx, addr := range serversAddress {
		fmt.Printf("%d %s\t", idx, addr)
	}
	fmt.Println()
	// fmt.Printf("status:\n%s\n", serversStatus)
	for idx, maxKey := range serversMaxKey {
		fmt.Printf("%d %d\t", idx, maxKey)
	}
	fmt.Println()

	// fmt.Printf("maxKey:\n%s\n", serversMaxKey)
	for idx, status := range serversStatus {
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
	addr := serversAddress[*shardIdx]
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
