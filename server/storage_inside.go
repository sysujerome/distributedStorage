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

	"github.com/syndtr/goleveldb/leveldb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	db         *common.SMap
	operations [][]string
	// serversAddress map[int64]string
	// serversStatus  map[int64]string
	// serversMaxKey  map[int64]int64
	serversAddress *common.IMap
	serversStatus  *common.IMap
	serversMaxKey  *common.IIMap
	shardIdx       = flag.Int64("shard_idx", 0, "the index of this shard node")
	next           = flag.Int64("next", 0, "the spilt node pointer")
	level          = flag.Int64("level", 0, "the initial level of hash function")
	hashSize       = flag.Int64("hash_size", 4, "the initial count of hash number")
	conf           = flag.String("conf", "conf.json", "the defaulted configure file")
	statusCode     common.Status
	errorCode      common.Error
	serverStatus   common.ServerStatus
	canSplit       bool
	overflowDB     *leveldb.DB
	spliting       int64
	version        int64
)

type server struct {
	pb.UnimplementedStorageServer
}

func (s *server) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {

	// 检查conf，判断该key是否为存在该节点上
	if hashFunc(in.GetKey()) != *shardIdx {
		target := hashFunc(in.GetKey())
		if target != *shardIdx {
			return &pb.GetReply{Result: "", Status: statusCode.Moved, Err: errorCode.Moved, Target: target, Version: version}, nil
		}
	}
	status, _ := serversStatus.ReadMap(*shardIdx)
	if status == serverStatus.Sleep {
		return &pb.GetReply{Result: "", Status: statusCode.Failed, Err: errorCode.NotWorking, Version: version}, nil
	}

	if atomic.LoadInt64(&spliting) > 0 {
		// operations = append(operations, []string{"get", in.GetKey()})
		// value, found := db[in.GetKey()]
		value, found := db.ReadMap(in.GetKey())
		if found {
			return &pb.GetReply{Result: value, Status: statusCode.Ok}, nil
		} else {
			// 从溢出表中查找
			// data, err := overflowDB.Get([]byte("key"), nil)
			// if (err != nil) {
			// 	return &pb.GetReply{Result: string(data), Status: statusCode.Ok}, nil
			// } else {
			secondIdx := int64(math.Pow(2, float64(*level)))*(*hashSize) + *shardIdx
			if int(secondIdx) >= serversAddress.GetLen() {
				return &pb.GetReply{Result: "", Status: statusCode.Failed, Err: errorCode.NotFound}, nil
			}
			address, _ := serversAddress.ReadMap(secondIdx)
			target := address
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
				return &pb.GetReply{Result: reply.GetResult(), Status: statusCode.Ok, Version: version}, nil
			} else {
				return &pb.GetReply{Result: "", Status: statusCode.Failed, Err: errorCode.NotFound, Version: version}, nil
			}
			// }
		}

		return &pb.GetReply{Result: "", Status: statusCode.Stored, Err: errorCode.Stored}, nil
	}
	// value, found := db[in.GetKey()]
	value, found := db.ReadMap(in.GetKey())
	if !found {
		// data, err := overflowDB.Get([]byte(in.GetKey()), nil)
		// if err != nil {
		// 	return &pb.GetReply{Result: "", Err: errorCode.NotFound, Status: statusCode.Failed}, nil
		// }
		// return &pb.GetReply{Result: string(data), Err: "", Status: statusCode.Ok}, nil

		// 存储在内部
		return &pb.GetReply{Result: "", Err: errorCode.NotFound, Status: statusCode.Failed}, nil
	}
	return &pb.GetReply{Result: value, Status: statusCode.Ok}, nil
}

func (s *server) Set(ctx context.Context, in *pb.SetRequest) (*pb.SetReply, error) {

	status, _ := serversStatus.ReadMap(*shardIdx)
	if status == serverStatus.Sleep {
		return &pb.SetReply{Result: "", Status: statusCode.Failed, Err: errorCode.NotWorking}, nil
	}
	// 检查conf，判断该key是否为存在该节点上
	if hashFunc(in.GetKey()) != *shardIdx {
		target := hashFunc(in.GetKey())
		if target != *shardIdx {
			return &pb.SetReply{Result: "", Status: statusCode.Moved, Err: errorCode.Moved, Target: target, Version: version}, nil
		}
	}
	if atomic.LoadInt64(&spliting) > 0 {
		// time.Sleep(time.Millisecond)
		// fmt.Printf("-")
		operations = append(operations, []string{"set", in.GetKey(), in.GetValue()})
		return &pb.SetReply{Result: "", Status: statusCode.Stored, Err: errorCode.Stored}, nil
	}
	// }

	// if len(db) == int(serversMaxKey[*shardIdx]) {
	maxKey, _ := serversMaxKey.ReadMap(*shardIdx)
	if db.GetLen() == int(maxKey) {
		// 继续放在本地db上
		// db[in.GetKey()] = in.GetValue()
		db.WriteMap(in.GetKey(), in.GetValue())
		if canSplit {
			split()
		}
	}
	// db[in.GetKey()] = in.GetValue()
	db.WriteMap(in.GetKey(), in.GetValue())
	value, found := db.ReadMap(in.GetKey())
	if !found || value != in.GetValue() {
		panic("error")
		return &pb.SetReply{Result: "", Status: statusCode.Failed, Err: errorCode.NotDefined}, nil
	}
	return &pb.SetReply{Status: statusCode.Ok}, nil
}

func (s *server) Del(ctx context.Context, in *pb.DelRequest) (*pb.DelReply, error) {
	status, _ := serversStatus.ReadMap(*shardIdx)
	if status == serverStatus.Sleep {
		return &pb.DelReply{Result: 0, Status: statusCode.Failed, Err: errorCode.NotWorking}, nil
	}

	// log.Printf("Received get key: %v", in.GetKey())
	// if serversStatus[*shardIdx] == serverStatus.Spliting {
	if atomic.LoadInt64(&spliting) > 0 {
		operations = append(operations, []string{"del", in.GetKey()})
		return &pb.DelReply{Result: 0, Status: statusCode.Stored, Err: errorCode.Stored}, nil
	}

	// delete(db, in.GetKey())
	db.Del(in.GetKey())

	_, found := db.ReadMap(in.GetKey())
	if found {
		return &pb.DelReply{Result: 0, Status: statusCode.Failed, Err: errorCode.NotDefined}, nil
	}
	return &pb.DelReply{Result: 1, Status: statusCode.Ok}, nil
}

func split() {
	if *next >= int64(serversAddress.GetLen()) {
		// fmt.Println("分裂失败, 节点满了。。。")
		return
	}
	addr, _ := serversAddress.ReadMap(*next)
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
		// fmt.Println("分裂失败, 节点满了。。。")
	}
}

func (s *server) Split(ctx context.Context, in *pb.SplitRequest) (*pb.SplitReply, error) {
	// log.Printf("Received get key: %v", in.GetKey())

	splitInside := func() { // 内部元素分裂
		// fmt.Printf("%s : 节点分裂中....\n", serversAddress[*shardIdx])
		secondIdx := int64(math.Pow(2, float64(*level)))**hashSize + *shardIdx
		if secondIdx >= int64(serversAddress.GetLen()) {
			// fmt.Printf("%s : 节点数满了....\n", serversAddress[*shardIdx])
			canSplit = false
			return
		}

		*next++
		if *next == int64(math.Pow(2, float64(*level)))**hashSize {
			*next = 0
			*level++
		}
		syncConf()
		target, found := serversAddress.ReadMap(secondIdx)
		judge(found == true)
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

		thisAddress, found := serversAddress.ReadMap(*shardIdx)
		judge(found == true)
		thisConn, err := grpc.Dial(thisAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		check(err)
		defer thisConn.Close()
		thisClient := pb.NewStorageClient(thisConn)
		// Contact the server and print out its response.

		// fmt.Printf("%s : 开始分裂....\n", serversAddress[secondIdx])
		count := 0

		// for key, value := range db.Map {
		key, value, exist := db.Range()
		for exist {

			// delete(db, key)
			db.Del(key)
			if hashFunc(key) == secondIdx {
				reply1, err := c.Set(ctx, &pb.SetRequest{Key: key, Value: value})
				check(err)
				if reply1.GetStatus() != statusCode.Ok && reply1.GetStatus() != statusCode.Stored {
					for reply1.GetStatus() == statusCode.Moved {
						target := reply1.GetTarget()
						targetAddress, _ := serversAddress.ReadMap(target)
						targetConn, err := grpc.Dial(targetAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
						check(err)
						defer targetConn.Close()
						targetClient := pb.NewStorageClient(targetConn)
						reply1, err = targetClient.Set(ctx, &pb.SetRequest{Key: key, Value: value})
						check(err)
					}
					if reply1.GetStatus() != statusCode.Ok && reply1.GetStatus() != statusCode.Stored {
						fmt.Println(reply1.GetStatus())
						fmt.Println("出错。。。")
						panic(reply1.GetErr())
					}
				}
				count++
			} else {
				reply2, err := thisClient.Set(ctx, &pb.SetRequest{Key: key, Value: value})
				check(err)
				if reply2.GetStatus() != statusCode.Ok && reply2.GetStatus() != statusCode.Stored {
					for reply2.GetStatus() == statusCode.Moved {
						target := reply2.GetTarget()
						targetAddress, found := serversAddress.ReadMap(target)
						judge(found == true)
						targetConn, err := grpc.Dial(targetAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
						check(err)
						defer targetConn.Close()
						targetClient := pb.NewStorageClient(targetConn)
						reply2, err = targetClient.Set(ctx, &pb.SetRequest{Key: key, Value: value})
						check(err)
					}
					if reply2.GetStatus() != statusCode.Ok && reply2.GetStatus() != statusCode.Stored {
						fmt.Println("出错。。。")
						panic(reply2.GetErr())
					}
				}
			}
			key, value, exist = db.Range()
		}
		// 遍历溢出表的内容
		// iter := overflowDB.NewIterator(nil, nil)
		// for iter.Next() {
		// 	key := string(iter.Key())
		// 	value := string(iter.Value())
		// 	curIDX := hashFunc(string(key))
		// 	err := overflowDB.Delete([]byte(key), nil)
		// 	check(err)

		// 	if curIDX == secondIdx {
		// 		reply1, err := c.Set(ctx, &pb.SetRequest{Key: key, Value: value})
		// 		check(err)
		// 		if reply1.GetStatus() != statusCode.Ok {
		// 			fmt.Println(reply1.GetStatus())
		// 			fmt.Println("出错。。。")
		// 			panic(reply1.GetErr())
		// 		}
		// 		count++
		// 	} else {
		// 		reply2, err := thisClient.Set(ctx, &pb.SetRequest{Key: key, Value: value})
		// 		check(err)
		// 		if reply2.GetStatus() != statusCode.Ok {
		// 			fmt.Println(reply2.GetStatus())
		// 			fmt.Println("出错。。。")
		// 			panic(reply2.GetErr())
		// 		}
		// 	}
		// }
		// serversStatus[*shardIdx] = serverStatus.Working
		serversStatus.WriteMap(*shardIdx, serverStatus.Working)
		syncConf()
		// fmt.Printf("%s : 分裂完成....\n", serversAddress[secondIdx])
		// 处理存起来的操作
		// fmt.Println("处理存起来的操作...")
		for _, operation := range operations {
			switch operation[0] {
			case "set":
				idx := hashFunc(operation[1])
				if idx == *shardIdx {
					// mutex.Lock()
					// db[operation[1]] = db[operation[2]]
					db.WriteMap(operation[1], operation[2])
					// mutex.Unlock()
				} else {
					target, found := serversAddress.ReadMap(idx)
					judge(found)
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
					// mutex.Lock()
					// db[operation[1]] = db[operation[2]]
					db.WriteMap(operation[1], operation[2])
					// mutex.Unlock()
				} else {
					target, found := serversAddress.ReadMap(idx)
					judge(found == true)
					conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
					check(err)
					defer conn.Close()
					c := pb.NewStorageClient(conn)
					c.Del(ctx, &pb.DelRequest{Key: operation[1]})
				}
			}
		}
		// serversStatus[*shardIdx] = serverStatus.Working
		serversStatus.WriteMap(*shardIdx, serverStatus.Working)
		return
	}

	// mutex.Lock()
	// serversStatus[*shardIdx] = serverStatus.Spliting
	serversStatus.WriteMap(*shardIdx, serverStatus.Spliting)

	atomic.AddInt64(&spliting, 1)
	// mutex.Unlock()
	go splitInside()

	// mutex.Lock()
	// serversStatus[*shardIdx] = serverStatus.Working
	serversStatus.WriteMap(*shardIdx, serverStatus.Working)
	atomic.AddInt64(&spliting, -1)
	version++
	syncConf()
	// mutex.Unlock()
	// if full {
	// 	canSplit = false
	// 	syncConf()
	// 	return &pb.SplitReply{Status: statusCode.Ok, Result: 0, Full: true}, nil
	// }

	return &pb.SplitReply{Status: statusCode.Ok, Version: version}, nil
}

func (s *server) Scan(ctx context.Context, in *pb.ScanRequest) (*pb.ScanReply, error) {
	// count := len(db)
	count := db.GetLen()
	address, found := serversAddress.ReadMap(*shardIdx)
	judge(found)
	result := address + "\n"
	result += fmt.Sprintf("%s", db.Map)
	return &pb.ScanReply{Result: result, Status: statusCode.Ok, Count: int64(count), Version: version}, nil
}

func syncConf() {
	request := pb.SyncConfRequest{}
	request.Begin = *shardIdx
	request.HashSize = *hashSize
	request.Level = *level
	request.Next = *next
	request.Server = make([]*pb.SyncConfRequest_ServConf, 0)
	request.CanSplit = canSplit
	idx, addr, exist := serversAddress.Range()
	for exist {
		var server pb.SyncConfRequest_ServConf
		server.Address = addr
		server.Idx = idx
		server.MaxKey, _ = serversMaxKey.ReadMap(idx)
		server.Status, _ = serversStatus.ReadMap(idx)
		request.Server = append(request.Server, &server)
		idx, addr, exist = serversAddress.Range()
	}
	target := *shardIdx + 1
	if target == int64(serversAddress.GetLen()) {
		target = 0
	}
	serv, found := serversAddress.ReadMap(target)
	judge(found)
	conn, err := grpc.Dial(serv, grpc.WithTransportCredentials(insecure.NewCredentials()))
	check(err)
	c := pb.NewStorageClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Second))
	defer cancel()
	reply, err := c.SyncConf(ctx, &request)
	check(err)
	if reply.GetStatus() != statusCode.Ok {
		address, _ := serversAddress.ReadMap(*shardIdx)
		fmt.Printf("%s : 同步失败。。。\n", address)
	}
}

func (s *server) SyncConf(ctx context.Context, in *pb.SyncConfRequest) (*pb.SyncConfReply, error) {
	if in.GetBegin() == *shardIdx {
		return &pb.SyncConfReply{Status: statusCode.Ok, Result: "同步结束", Version: version}, nil
	}
	// fmt.Printf("%s : 节点同步中...\n", serversAddress[*shardIdx])
	*next = in.GetNext()
	*level = in.GetLevel()
	*hashSize = in.GetHashSize()
	canSplit = in.GetCanSplit()
	serversAddress = &common.IMap{
		Map: make(map[int64]string),
	}
	serversStatus = &common.IMap{
		Map: make(map[int64]string),
	}
	serversMaxKey = &common.IIMap{
		Map: make(map[int64]int64),
	}
	for _, serv := range in.GetServer() {
		serversAddress.WriteMap(serv.Idx, serv.Address)
		serversMaxKey.WriteMap(serv.Idx, serv.MaxKey)
		serversStatus.WriteMap(serv.Idx, serv.Status)
	}
	success := false
	target := *shardIdx + 1
	if target == int64(serversAddress.GetLen()) {
		target = 0
	}
	for !success {
		serv, found := serversAddress.ReadMap(target)
		judge(found == true)
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
		if target == int64(serversAddress.GetLen()) {
			target = 0
		}
	}
	return &pb.SyncConfReply{Status: statusCode.Ok, Version: version}, nil
}

func (s *server) GetConf(ctx context.Context, in *pb.GetConfRequest) (*pb.GetConfReply, error) {
	var reply pb.GetConfReply
	reply.HashSize = *hashSize
	reply.Level = *level
	reply.Next = *next
	reply.Status = statusCode.Ok
	idx, addr, exist := serversAddress.Range()
	for exist {
		var server pb.GetConfReply_ServConf
		server.Address = addr
		server.Idx = idx
		server.MaxKey, _ = serversMaxKey.ReadMap(idx)
		server.Status, _ = serversStatus.ReadMap(idx)
		reply.Server = append(reply.Server, &server)
		idx, addr, exist = serversAddress.Range()
	}
	reply.Version = version
	return &reply, nil
}

func (s *server) WakeUp(ctx context.Context, in *pb.WakeRequest) (*pb.WakeReply, error) {
	status, found := serversStatus.ReadMap(*shardIdx)
	judge(found == true)
	if status != serverStatus.Sleep {
		return &pb.WakeReply{Status: statusCode.Failed, Result: "该节点不在休眠状态", Version: version}, nil
	}

	serversStatus.WriteMap(*shardIdx, serverStatus.Working)
	syncConf()
	return &pb.WakeReply{Status: statusCode.Ok, Version: version}, nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func judge(b bool) {
	if b == false {
		panic("something wrong")
	}
}

func initConf() {
	flag.Parse()

	// db = make(map[string]string)
	db = &common.SMap{
		Map: make(map[string]string),
	}
	overflowDB, _ = leveldb.OpenFile("./overflow.db", nil)
	// check(err)
	operations = make([][]string, 0)
	// serversAddress = make(map[int64]string)
	// serversStatus = make(map[int64]string)
	// serversMaxKey = make(map[int64]int64)
	serversAddress = &common.IMap{
		Map: make(map[int64]string),
	}
	serversStatus = &common.IMap{
		Map: make(map[int64]string),
	}
	serversMaxKey = &common.IIMap{
		Map: make(map[int64]int64),
	}
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
		serversAddress.WriteMap(idx, address)
		serversStatus.WriteMap(idx, ShardNodeConfs["status"].(string))
		serversMaxKey.WriteMap(idx, int64(ShardNodeConfs["max_key"].(float64)))
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
	fmt.Printf("db:\n%s\n", db.Map)
	// fmt.Printf("servers:\n%s\n", serversAddress)
	idx, addr, exist := serversAddress.Range()
	for exist {
		fmt.Printf("%d %s\t", idx, addr)
		idx, addr, exist = serversAddress.Range()
	}
	fmt.Println()
	// fmt.Printf("status:\n%s\n", serversStatus)
	idx, maxKey, exist := serversMaxKey.Range()
	for exist {
		fmt.Printf("%d %d\t", idx, maxKey)
		idx, maxKey, exist = serversMaxKey.Range()
	}
	fmt.Println()

	// fmt.Printf("maxKey:\n%s\n", serversMaxKey)
	idx, status, exist := serversStatus.Range()
	for exist {
		fmt.Printf("%d %s\t", idx, status)
		idx, status, exist = serversStatus.Range()
	}
	fmt.Println()
	fmt.Printf("shardIdx:\n%d\n", *shardIdx)
	fmt.Printf("next:\n%d\n", *next)
	fmt.Printf("level:\n%d\n", *level)
	fmt.Printf("hashSize:\n%d\n", *hashSize)
	fmt.Printf("conf:\n%s\n", *conf)
}

func serve() {
	addr, found := serversAddress.ReadMap(*shardIdx)
	judge(found == true)
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
