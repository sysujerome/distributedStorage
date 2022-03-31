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
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// 集群的元数据
type Meta struct {
	sync.RWMutex   //避免对元数据的竞态访问
	serversAddress map[int64]string
	serversStatus  map[int64]string
	serversMaxKey  map[int64]int64
	next           int64
	level          int64
	hashSize       int64
	version        int64
}

var (
	meta         Meta //集群的元数据
	db           DataBase
	operations   [][]string
	shardIdx     = flag.Int64("shard_idx", 0, "the index of this shard node")
	conf         = flag.String("conf", "conf.json", "the defaulted configure file")
	statusCode   common.Status
	errorCode    common.Error
	serverStatus common.ServerStatus
	canSplit     bool
	// mutex          sync.Mutex
	spliting  int64
	confMutex sync.RWMutex
)

type server struct {
	pb.UnimplementedStorageServer
}

func (s *server) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {

	meta.RLock()
	meta0 := meta
	meta.RUnlock()
	// 检查conf，判断该key是否为存在该节点上
	if hashFunc(in.GetKey()) != *shardIdx {
		target := hashFunc(in.GetKey())
		if target != *shardIdx {
			return &pb.GetReply{Result: "", Status: statusCode.Moved, Err: errorCode.Moved, Target: target, Version: meta0.version}, nil
		}
	}

	if meta0.serversStatus[*shardIdx] == serverStatus.Sleep {
		return &pb.GetReply{Result: "", Status: statusCode.Failed, Err: errorCode.NotWorking, Version: meta0.version}, nil
	}
	if atomic.LoadInt64(&spliting) > 0 {
		// operations = append(operations, []string{"get", in.GetKey()})
		found := db.get(in.GetKey())
		if found {
			return &pb.GetReply{Result: "", Status: statusCode.Ok}, nil
		} else {
			// 从溢出表中查找
			// data, err := overflowDB.Get([]byte("key"), nil)
			// if (err != nil) {
			// 	return &pb.GetReply{Result: string(data), Status: statusCode.Ok}, nil
			// } else {
			secondIdx := int64(math.Pow(2, float64(meta0.level)))*(meta0.hashSize) + *shardIdx
			if int(secondIdx) >= len(meta0.serversAddress) {
				return &pb.GetReply{Result: "", Status: statusCode.Failed, Err: errorCode.NotFound}, nil
			}
			target := meta0.serversAddress[secondIdx]
			conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
			check(err)
			defer conn.Close()
			c := pb.NewStorageClient(conn)
			// Contact the server and print out its response.
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(2*time.Second))
			defer cancel()

			reply, err := c.Get(ctx, &pb.GetRequest{Key: in.GetKey()})
			check(err)
			if reply.GetStatus() == statusCode.Ok {
				return &pb.GetReply{Result: reply.GetResult(), Status: statusCode.Ok, Version: meta0.version}, nil
			} else {
				return &pb.GetReply{Result: "", Status: statusCode.Failed, Err: errorCode.NotFound, Version: meta0.version}, nil
			}
			// }
		}

		return &pb.GetReply{Result: "", Status: statusCode.Stored, Err: errorCode.Stored}, nil
	}
	found := db.get(in.GetKey())
	if !found {
		// data, err := overflowDB.Get([]byte(in.GetKey()), nil)
		// if err != nil {
		// 	return &pb.GetReply{Result: "", Err: errorCode.NotFound, Status: statusCode.Failed}, nil
		// }
		// return &pb.GetReply{Result: string(data), Err: "", Status: statusCode.Ok}, nil

		// 存储在内部
		return &pb.GetReply{Result: "", Err: errorCode.NotFound, Status: statusCode.Failed}, nil
	}
	return &pb.GetReply{Result: "", Status: statusCode.Ok}, nil
}

func (s *server) Set(ctx context.Context, in *pb.SetRequest) (*pb.SetReply, error) {
	meta.RLock()
	meta0 := meta
	meta.RUnlock()
	if meta0.serversStatus[*shardIdx] == serverStatus.Sleep {
		return &pb.SetReply{Result: "", Status: statusCode.Failed, Err: errorCode.NotWorking}, nil
	}
	// 检查conf，判断该key是否为存在该节点上
	if hashFunc(in.GetKey()) != *shardIdx {
		target := hashFunc(in.GetKey())
		if target != *shardIdx {
			return &pb.SetReply{Result: "", Status: statusCode.Moved, Err: errorCode.Moved, Target: target, Version: meta0.version}, nil
		}
	}
	if atomic.LoadInt64(&spliting) > 0 {
		// time.Sleep(time.Millisecond)
		// fmt.Printf("-")
		operations = append(operations, []string{"set", in.GetKey()})
		return &pb.SetReply{Result: "", Status: statusCode.Stored, Err: errorCode.Stored}, nil
	}
	// }

	if int64(db.len) == meta0.serversMaxKey[*shardIdx] {
		// 继续放在本地db上
		// db[in.GetKey()] = in.GetValue()
		db.set(in.GetKey())
		if canSplit {
			split()
		}
	}
	db.set(in.GetKey())
	return &pb.SetReply{Status: statusCode.Ok}, nil
}

func (s *server) Del(ctx context.Context, in *pb.DelRequest) (*pb.DelReply, error) {
	meta.RLock()
	meta0 := meta
	meta.RUnlock()
	if meta0.serversStatus[*shardIdx] == serverStatus.Sleep {
		return &pb.DelReply{Result: 0, Status: statusCode.Failed, Err: errorCode.NotWorking}, nil
	}

	// log.Printf("Received get key: %v", in.GetKey())
	// if serversStatus[*shardIdx] == serverStatus.Spliting {
	if atomic.LoadInt64(&spliting) > 0 {
		operations = append(operations, []string{"del", in.GetKey()})
		return &pb.DelReply{Result: 0, Status: statusCode.Stored, Err: errorCode.Stored}, nil
	}
	// delete(db, in.GetKey())
	success := db.del(in.GetKey())
	if success {
		return &pb.DelReply{Result: 0, Status: statusCode.Failed, Err: errorCode.NotDefined}, nil
	}
	return &pb.DelReply{Result: 1, Status: statusCode.Ok}, nil
}

func split() {
	meta.RLock()
	meta0 := meta
	meta.RUnlock()
	if meta0.next >= int64(len(meta0.serversAddress)) {
		// fmt.Println("分裂失败, 节点满了。。。")
		return
	}
	addr := meta0.serversAddress[meta0.next]
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	check(err)
	defer conn.Close()
	c := pb.NewStorageClient(conn)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Second))
	defer cancel()

	reply, err := c.Split(ctx, &pb.SplitRequest{})
	check(err)
	if reply.GetFull() {
		// fmt.Println("分裂失败, 节点满了。。。")
	}
}

func (s *server) Split(ctx context.Context, in *pb.SplitRequest) (*pb.SplitReply, error) {
	// log.Printf("Received get key: %v", in.GetKey())

	meta.RLock()
	meta0 := meta
	meta.RUnlock()
	// 分裂时一直新建连接会出现too many open files错误
	var clients []pb.StorageClient
	for idx := 0; idx < len(meta0.serversAddress); idx++ {
		addr := meta0.serversAddress[int64(idx)]
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		check(err)
		c := pb.NewStorageClient(conn)
		clients = append(clients, c)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Second))
	defer cancel()

	splitInside := func() { // 内部元素分裂
		fmt.Printf("%s:  %s : 节点分裂中....\n", time.Now().Format("2006-01-02 15:04:05"), meta0.serversAddress[*shardIdx])
		secondIdx := int64(math.Pow(2, float64(meta0.level)))*meta0.hashSize + *shardIdx
		if secondIdx >= int64(len(meta0.serversAddress)) {
			// fmt.Printf("%s : 节点数满了....\n", serversAddress[*shardIdx])
			canSplit = false
			return
		}
		meta.Lock()
		meta0.next++
		if meta0.next == int64(math.Pow(2, float64(meta0.level)))*meta0.hashSize {
			meta0.next = 0
			meta0.level++
		}
		meta0.next = meta.next
		meta.Unlock()
		syncConf()
		// target := serversAddress[secondIdx]
		// conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
		// check(err)
		// defer conn.Close()
		// c := pb.NewStorageClient(conn)
		// Contact the server and print out its response.

		reply, err := clients[secondIdx].WakeUp(ctx, &pb.WakeRequest{})
		check(err)
		if reply.GetStatus() != statusCode.Ok {
			// fmt.Printf("%s %s : 节点唤醒失败....\n", serversAddress[*shardIdx], serversAddress[secondIdx])
			return
		}
		syncConf()

		// thisAddress := serversAddress[]
		// thisConn, err := grpc.Dial(thisAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		// check(err)
		// defer thisConn.Close()
		// thisClient := pb.NewStorageClient(thisConn)
		// Contact the server and print out its response.

		// fmt.Printf("%s : 开始分裂....\n", serversAddress[secondIdx])
		count := 0
		// for key, value := range db {
		// 	delete(db, key)
		// 	if hashFunc(key) == secondIdx {
		// 		reply1, err := clients[secondIdx].Set(ctx, &pb.SetRequest{Key: key, Value: value})
		// 		check(err)
		// 		if reply1.GetStatus() != statusCode.Ok && reply1.GetStatus() != statusCode.Stored {
		// 			for reply1.GetStatus() == statusCode.Moved {
		// 				reply1, err = clients[int(reply1.GetTarget())].Set(ctx, &pb.SetRequest{Key: key, Value: value})
		// 				check(err)
		// 			}
		// 			if reply1.GetStatus() != statusCode.Ok && reply1.GetStatus() != statusCode.Stored {
		// 				fmt.Println(reply1.GetStatus())
		// 				fmt.Println("出错。。。")
		// 				panic(reply1.GetErr())
		// 			}
		// 		}
		// 		count++
		// 	} else {
		// 		reply2, err := clients[*shardIdx].Set(ctx, &pb.SetRequest{Key: key, Value: value})
		// 		check(err)
		// 		if reply2.GetStatus() != statusCode.Ok && reply2.GetStatus() != statusCode.Stored {
		// 			for reply2.GetStatus() == statusCode.Moved {
		// 				reply2, err = clients[int(reply2.GetTarget())].Set(ctx, &pb.SetRequest{Key: key, Value: value})
		// 				check(err)
		// 			}
		// 			if reply2.GetStatus() != statusCode.Ok && reply2.GetStatus() != statusCode.Stored {
		// 				fmt.Println(reply2.GetStatus())
		// 				fmt.Println("出错。。。")
		// 				panic(reply2.GetErr())
		// 			}
		// 		}
		// 	}
		// }

		for index := 0; index < length; index++ {
			db.slots[index].Lock()
			tempData := make([]string, 0) //学习redis来进行渐进性rehash
			for _, key := range db.slots[index].data {
				if hashFunc(key) == secondIdx {
					reply1, err := clients[secondIdx].Set(ctx, &pb.SetRequest{Key: key})
					check(err)
					if reply1.GetStatus() != statusCode.Ok && reply1.GetStatus() != statusCode.Stored {
						for reply1.GetStatus() == statusCode.Moved {
							reply1, err = clients[int(reply1.GetTarget())].Set(ctx, &pb.SetRequest{Key: key})
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
					tempData = append(tempData, key)
				}
			}
			db.slots[index].data = tempData
			db.slots[index].Unlock()
		}

		meta.Lock()
		meta.serversStatus[*shardIdx] = serverStatus.Working
		meta.Unlock()
		syncConf()
		// fmt.Printf("%s : 分裂完成....\n", serversAddress[secondIdx])
		// 处理存起来的操作
		// fmt.Println("处理存起来的操作...")
		for _, operation := range operations {
			switch operation[0] {
			case "set":
				idx := hashFunc(operation[1])
				if idx == *shardIdx {
					// db[operation[1]] = db[operation[2]]
					db.set(operation[1])
				} else {
					reply, err := clients[idx].Set(ctx, &pb.SetRequest{Key: operation[1]})
					if err != nil {
						panic(err)
					}
					if reply.GetStatus() != statusCode.Ok {
						panic(reply.GetStatus())
					}
				}
				// db[operation[1]] = operation[2]
			case "del":
				// delete(db, operation[1])
				idx := hashFunc(operation[1])
				if idx == *shardIdx {
					// db[operation[1]] = db[operation[2]]
					db.del(operation[1])
				} else {
					clients[idx].Del(ctx, &pb.DelRequest{Key: operation[1]})
				}
			}
		}
		meta.Lock()
		meta.serversStatus[*shardIdx] = serverStatus.Working
		meta.Unlock()
	}

	meta.Lock()
	meta.serversStatus[*shardIdx] = serverStatus.Spliting
	meta.Unlock()
	atomic.AddInt64(&spliting, 1)
	// mutex.Unlock()
	splitInside()

	// mutex.Lock()
	meta.Lock()
	meta.serversStatus[*shardIdx] = serverStatus.Working
	atomic.AddInt64(&spliting, -1)
	meta.version++
	meta.Unlock()
	meta0.version++
	syncConf()
	// mutex.Unlock()
	// if full {
	// 	canSplit = false
	// 	syncConf()
	// 	return &pb.SplitReply{Status: statusCode.Ok, Result: 0, Full: true}, nil
	// }

	return &pb.SplitReply{Status: statusCode.Ok, Version: meta0.version}, nil
}

func (s *server) Scan(ctx context.Context, in *pb.ScanRequest) (*pb.ScanReply, error) {

	meta.RLock()
	meta0 := meta
	meta.RUnlock()
	count := db.len
	result := meta0.serversAddress[*shardIdx] + "\n"
	// result += fmt.Sprintf("%s", db.scan())
	result += db.scan()
	return &pb.ScanReply{Result: result, Status: statusCode.Ok, Count: int64(count), Version: meta0.version}, nil
}

func syncConf() {
	meta.RLock()
	meta0 := meta
	meta.RUnlock()
	request := pb.SyncConfRequest{}
	request.Begin = *shardIdx
	request.HashSize = meta0.hashSize
	request.Level = meta0.level
	request.Next = meta0.next
	request.Server = make([]*pb.SyncConfRequest_ServConf, 0)
	request.CanSplit = canSplit
	for idx, addr := range meta0.serversAddress {
		var server pb.SyncConfRequest_ServConf
		server.Address = addr
		server.Idx = idx
		server.MaxKey = meta0.serversMaxKey[idx]
		server.Status = meta0.serversStatus[idx]
		request.Server = append(request.Server, &server)
	}
	target := *shardIdx + 1
	if target == int64(len(meta0.serversAddress)) {
		target = 0
	}
	serv := meta0.serversAddress[target]
	conn, err := grpc.Dial(serv, grpc.WithTransportCredentials(insecure.NewCredentials()))
	check(err)
	c := pb.NewStorageClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(2*time.Second))
	defer cancel()
	reply, err := c.SyncConf(ctx, &request)
	check(err)
	if reply.GetStatus() != statusCode.Ok {
		fmt.Printf("%s : 同步失败。。。\n", meta0.serversAddress[*shardIdx])
	}
}

func (s *server) SyncConf(ctx context.Context, in *pb.SyncConfRequest) (*pb.SyncConfReply, error) {
	meta.RLock()
	version := meta.version
	meta.RUnlock()
	if in.GetBegin() == *shardIdx {
		return &pb.SyncConfReply{Status: statusCode.Ok, Result: "同步结束", Version: version}, nil
	}
	// 更新集群元数据
	meta.Lock()
	meta.next = in.GetNext()
	meta.level = in.GetLevel()
	meta.hashSize = in.GetHashSize()

	// 为啥要先删除呢？
	// for k := range serversAddress {
	// 	delete(serversAddress, k)
	// 	delete(serversMaxKey, k)
	// 	delete(serversStatus, k)
	// }
	for _, serv := range in.GetServer() {
		meta.serversAddress[serv.Idx] = serv.Address
		meta.serversMaxKey[serv.Idx] = serv.MaxKey
		meta.serversStatus[serv.Idx] = serv.Status
	}
	meta.Unlock()
	// 线性传递下去
	// 只要能将meta数据正确传递给下一个节点就行，如果id+1的节点不通就id+2...
	success := false
	target := *shardIdx + 1

	meta.RLock()
	meta0 := meta
	meta.RUnlock()
	if target == int64(len(meta0.serversAddress)) {
		target = 0
	}
	for !success {
		serv := meta0.serversAddress[target]
		conn, err := grpc.Dial(serv, grpc.WithTransportCredentials(insecure.NewCredentials()))
		check(err)
		c := pb.NewStorageClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(2*time.Second))
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
		if target == int64(len(meta0.serversAddress)) {
			target = 0
		}
	}
	return &pb.SyncConfReply{Status: statusCode.Ok, Version: version}, nil
}

func (s *server) GetConf(ctx context.Context, in *pb.GetConfRequest) (*pb.GetConfReply, error) {

	meta.RLock()
	meta0 := meta
	meta.RUnlock()
	var reply pb.GetConfReply
	reply.HashSize = meta0.hashSize
	reply.Level = meta0.level
	reply.Next = meta0.next
	reply.Status = statusCode.Ok
	for idx, addr := range meta0.serversAddress {
		var server pb.GetConfReply_ServConf
		server.Address = addr
		server.Idx = idx
		server.MaxKey = meta0.serversMaxKey[idx]
		server.Status = meta0.serversStatus[idx]
		reply.Server = append(reply.Server, &server)
	}
	reply.Version = meta0.version
	return &reply, nil
}

func (s *server) WakeUp(ctx context.Context, in *pb.WakeRequest) (*pb.WakeReply, error) {
	meta.RLock()
	meta0 := meta
	meta.RUnlock()
	fmt.Printf("%s : %s 被唤醒\n", time.Now().Format("2006-01-02 15:04:05"), meta0.serversAddress[*shardIdx])
	if meta0.serversStatus[*shardIdx] != serverStatus.Sleep {
		return &pb.WakeReply{Status: statusCode.Failed, Result: "该节点不在休眠状态", Version: meta0.version}, nil
	}
	meta.Lock()
	meta.serversStatus[*shardIdx] = serverStatus.Working
	meta.version++
	meta.Unlock()
	meta0.version++
	syncConf()
	return &pb.WakeReply{Status: statusCode.Ok, Version: meta0.version}, nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func initConf() {
	flag.Parse()

	// check(err)
	operations = make([][]string, 0)
	meta.Lock()
	meta.hashSize = 4
	meta.serversAddress = make(map[int64]string)
	meta.serversStatus = make(map[int64]string)
	meta.serversMaxKey = make(map[int64]int64)
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
		meta.serversAddress[idx] = address
		meta.serversStatus[idx] = ShardNodeConfs["status"].(string)
		meta.serversMaxKey[idx] = int64(ShardNodeConfs["max_key"].(float64))
		// fmt.Println(idx, ip, port)
	}
	meta.level = 0
	meta.Unlock()
	statusCode.Init()
	errorCode.Init()
	serverStatus.Init()
}

func hashFunc(key string) int64 {
	meta.RLock()
	meta0 := meta
	meta.RUnlock()
	posCRC16 := int64(crc16.Checksum([]byte(key), crc16.IBMTable))
	pos := posCRC16 % (int64(math.Pow(2.0, float64(meta0.level))) * meta0.hashSize)
	if pos < meta0.next { // 分裂过了的
		pos = posCRC16 % (int64(math.Pow(2.0, float64(meta0.level+1))) * meta0.hashSize)
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
	meta.RLock()
	meta0 := meta
	meta.RUnlock()
	fmt.Printf("db:\n%s\n", db.scan())
	// fmt.Printf("servers:\n%s\n", serversAddress)
	for idx, addr := range meta0.serversAddress {
		fmt.Printf("%d %s\t", idx, addr)
	}
	fmt.Println()
	// fmt.Printf("status:\n%s\n", serversStatus)
	for idx, maxKey := range meta0.serversMaxKey {
		fmt.Printf("%d %d\t", idx, maxKey)
	}
	fmt.Println()

	// fmt.Printf("maxKey:\n%s\n", serversMaxKey)
	for idx, status := range meta0.serversStatus {
		fmt.Printf("%d %s\t", idx, status)
	}
	fmt.Println()
	fmt.Printf("shardIdx:\n%d\n", *shardIdx)
	fmt.Printf("next:\n%d\n", meta0.next)
	fmt.Printf("level:\n%d\n", meta0.level)
	fmt.Printf("hashSize:\n%d\n", meta0.hashSize)
	fmt.Printf("conf:\n%s\n", *conf)
}

func serve() {
	meta.RLock()
	addr := meta.serversAddress[*shardIdx]
	meta.RUnlock()
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
