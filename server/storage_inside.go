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
	spliting int64
	maxKey   int64
	status   int32 // 0 for work, 1 for sleep, 2 for splitting, 3 for full
)

type server struct {
	pb.UnimplementedStorageServer
}

// // 每个增删改查操作都先检查当前节点状态
// // bool : 0 -> sleep, 1 -> work
// // int64 :
// func checkStatus() (bool, int64) {

// }

func (s *server) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
	// 先检查当前节点状态
	if atomic.LoadInt32(&status) == 1 {
		version := atomic.LoadInt64(&meta.version)
		return &pb.GetReply{Result: "", Status: statusCode.Failed, Err: errorCode.NotWorking, Version: version}, nil
	}
	version := atomic.LoadInt64(&meta.version)
	// 检查conf，判断该key是否为存在该节点上
	if hashFunc(in.GetKey()) != *shardIdx {
		target := hashFunc(in.GetKey())
		if target != *shardIdx {
			return &pb.GetReply{Result: "", Status: statusCode.Moved, Err: errorCode.Moved, Target: target, Version: version}, nil
		}
	}

	if atomic.LoadInt32(&status) == 2 {
		// operations = append(operations, []string{"get", in.GetKey()})
		found := db.get(in.GetKey())
		if found {
			return &pb.GetReply{Result: "", Status: statusCode.Ok}, nil
		} else {
			addresses := make(map[int64]string)
			meta.Lock()
			for k, v := range meta.serversAddress {
				addresses[k] = v
			}
			meta.Unlock()
			secondIdx := int64(math.Pow(2, float64(atomic.LoadInt64(&meta.level))))*(atomic.LoadInt64(&meta.hashSize)) + *shardIdx
			if int(secondIdx) >= len(addresses) {
				return &pb.GetReply{Result: "", Status: statusCode.Failed, Err: errorCode.NotFound}, nil
			}
			target := addresses[secondIdx]
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
				return &pb.GetReply{Result: reply.GetResult(), Status: statusCode.Ok, Version: version}, nil
			} else {
				return &pb.GetReply{Result: "", Status: statusCode.Failed, Err: errorCode.NotFound, Version: version}, nil
			}
			// }
		}
	}
	found := db.get(in.GetKey())
	if !found {
		// data, err := overflowDB.Get([]byte(in.GetKey()), nil)
		// if err != nil {
		// 	return &pb.GetReply{Result: "", Err: errorCode.NotFound, Status: statusCode.Failed}, nil
		// }
		// return &pb.GetReply{Result: string(data), Err: "", Status: statusCode.Ok}, nil

		// 存储在内部
		return &pb.GetReply{Result: "", Err: errorCode.NotFound, Status: statusCode.Failed, Version: version}, nil
	}
	return &pb.GetReply{Result: "", Status: statusCode.Ok, Version: version}, nil
}

func (s *server) Set(ctx context.Context, in *pb.SetRequest) (*pb.SetReply, error) {
	// 先检查当前节点状态
	// 检查当前节点状态为1，此节点正在休眠，返回fail和notworking的错误码
	if atomic.LoadInt32(&status) == 1 {
		version := atomic.LoadInt64(&meta.version)
		return &pb.SetReply{Result: "", Status: statusCode.Failed, Err: errorCode.NotWorking, Version: version}, nil
	}
	// 检查conf，判断该key是否为存在该节点上
	if hashFunc(in.GetKey()) != *shardIdx {
		target := hashFunc(in.GetKey())

		if target != *shardIdx {
			return &pb.SetReply{Result: "", Status: statusCode.Moved,
				Err: errorCode.Moved, Target: target, Version: atomic.LoadInt64(&meta.version),
				Next: atomic.LoadInt64(&meta.next), Level: atomic.LoadInt64(&meta.level)}, nil
		}
	}
	// 节点正在分裂
	if atomic.LoadInt32(&status) == 2 {
		// time.Sleep(time.Millisecond)
		// fmt.Printf("-")
		operations = append(operations, []string{"set", in.GetKey()})
		fmt.Printf("Store : %s\n", in.GetKey())
		return &pb.SetReply{Result: "", Status: statusCode.Stored, Err: errorCode.Stored}, nil
	}
	// }

	if int64(db.len) >= maxKey {
		// 继续放在本地db上
		// 将maxKey值扩大为1.25倍
		maxKey = int64(float64(maxKey) * 1.25)
		db.set(in.GetKey())
		if canSplit {
			split()
		}
	} else {
		db.set(in.GetKey())
	}
	return &pb.SetReply{Status: statusCode.Ok}, nil
}

func (s *server) Del(ctx context.Context, in *pb.DelRequest) (*pb.DelReply, error) {
	// 先检查当前节点状态
	if atomic.LoadInt32(&status) == 1 {
		version := atomic.LoadInt64(&meta.version)
		return &pb.DelReply{Result: 0, Status: statusCode.Failed, Err: errorCode.NotWorking, Version: version}, nil
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
	addresses := make(map[int64]string)
	meta.Lock()
	for k, v := range meta.serversAddress {
		addresses[k] = v
	}
	meta.Unlock()
	if atomic.LoadInt64(&meta.next) >= int64(len(addresses)) {
		// fmt.Println("分裂失败, 节点满了。。。")
		return
	}
	addr := addresses[atomic.LoadInt64(&meta.next)]
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	check(err)
	defer conn.Close()
	c := pb.NewStorageClient(conn)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Second))
	defer cancel()

	reply, err := c.Split(ctx, &pb.SplitRequest{Next: atomic.LoadInt64(&meta.next), Level: atomic.LoadInt64(&meta.level)})
	check(err)
	retry := 0
	for reply.GetStatus() == statusCode.Moved {
		retry++
		if retry > 16 {
			panic("Too much retry... in split moved.")
		}
		addr = addresses[reply.GetNext()]
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		check(err)
		defer conn.Close()
		c = pb.NewStorageClient(conn)
		// Contact the server and print out its response.
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Second))
		defer cancel()
		reply, err = c.Split(ctx, &pb.SplitRequest{Next: atomic.LoadInt64(&meta.next), Level: atomic.LoadInt64(&meta.level)})
		check(err)
	}
	if reply.GetStatus() != statusCode.Ok {
		fmt.Printf("server split : %s\n", reply.GetStatus())
		panic(fmt.Sprintf("server split failed : %s", addr))
	}
}

func (s *server) Split(ctx context.Context, in *pb.SplitRequest) (*pb.SplitReply, error) {
	// log.Printf("Received get key: %v", in.GetKey())
	// 先检查当前节点状态
	if atomic.LoadInt32(&status) == 1 {
		version := atomic.LoadInt64(&meta.version)
		return &pb.SplitReply{Result: 0, Status: statusCode.Failed, Err: errorCode.NotWorking, Version: version}, nil
	}
	// 排除并发分裂
	if atomic.LoadInt32(&status) == 2 {
		version := atomic.LoadInt64(&meta.version)
		return &pb.SplitReply{Status: statusCode.Ok, Version: version}, nil
	}

	// 检查集群元数据是否一致
	// 被分裂过的节点的next,level值相对更新，则通知去新的next节点分裂
	if atomic.LoadInt64(&meta.next) > in.GetNext() && atomic.LoadInt64(&meta.level) >= in.GetNext() || (atomic.LoadInt64(&meta.level) > in.GetLevel()) {
		return &pb.SplitReply{Result: 0, Status: statusCode.Moved, Next: atomic.LoadInt64(&meta.next), Level: atomic.LoadInt64(&meta.level)}, nil
	}

	// 开始分裂
	// 首先要保证集群的元数据一致性
	// atomic.SwapInt64(&meta.next, in.GetNext())
	// atomic.SwapInt64(&meta.level, in.GetLevel())

	// 当前节点的meta.next必须和当前的shardIdx一致
	idx := atomic.LoadInt64(shardIdx)
	atomic.SwapInt64(&meta.next, idx)

	addresses := make(map[int64]string)
	meta.RLock()
	for k, v := range meta.serversAddress {
		addresses[k] = v
	}
	meta.RUnlock()

	// 分裂时一直新建连接会出现too many open files错误
	// 和每个节点建立一个长连接
	var clients []pb.StorageClient
	for idx := 0; idx < len(addresses); idx++ {
		addr := addresses[int64(idx)]
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		check(err)
		c := pb.NewStorageClient(conn)
		clients = append(clients, c)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Second))
	defer cancel()

	splitInside := func() { // 内部元素分裂
		// 原子比较该节点是不是正在分裂，保证只有该节点只分裂一次
		if !atomic.CompareAndSwapInt32(&status, 0, 2) {
			return
		}
		atomic.StoreInt32(&status, 2)
		defer atomic.StoreInt32(&status, 0)
		fmt.Printf("节点分裂中...%s:  %s\n", addresses[*shardIdx], time.Now().Format("2006-01-02 15:04:05"))
		secondIdx := int64(math.Pow(2, float64(atomic.LoadInt64(&meta.level))))*atomic.LoadInt64(&meta.hashSize) + *shardIdx
		if secondIdx >= int64(len(addresses)) {
			fmt.Printf("%s : 节点数满了....\n", addresses[*shardIdx])
			canSplit = false
			return
		}

		// 分裂完成，后置条件：将next指针往后移动
		atomic.AddInt64(&meta.next, 1)
		// 一轮分裂已经完成，将next指针归零
		if atomic.LoadInt64(&meta.next) == int64(math.Pow(2, float64(atomic.LoadInt64(&meta.level))))*atomic.LoadInt64(&meta.hashSize) {
			atomic.StoreInt64(&meta.next, 0)
		}
		// 该节点已经分裂过，所以针对该节点的数据，是需要将level值加1后重新计算的
		atomic.AddInt64(&meta.level, 1)
		fmt.Printf("Next : %d\tLevel : %d\n", atomic.LoadInt64(&meta.next), atomic.LoadInt64(&meta.level))

		// 需要将next指针改变之后的数据同步到对应的wakeup节点，不然会出现同一个key的元数据不一致报错  (0->4 ==> 4->0 死循环)
		reply, err := clients[secondIdx].WakeUp(ctx, &pb.WakeRequest{Next: atomic.LoadInt64(&meta.next), Level: atomic.LoadInt64(&meta.level)})
		check(err)
		if reply.GetStatus() != statusCode.Ok {
			fmt.Printf("节点唤醒失败 : %s %s ....\n", addresses[*shardIdx], addresses[secondIdx])
			log.Panicf("err : %s\n", reply.GetErr())
			return
		}
		fmt.Printf("节点唤醒成功: %s %s ...\n", addresses[*shardIdx], addresses[secondIdx])
		// syncConf()

		// thisAddress := serversAddress[]
		// thisConn, err := grpc.Dial(thisAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		// check(err)
		// defer thisConn.Close()
		// thisClient := pb.NewStorageClient(thisConn)
		// Contact the server and print out its response.

		// fmt.Printf("%d : 开始分裂....\n", shardIdx)
		count := 0
		set := func(idx int64, key string) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(100*time.Second))
			defer cancel()
			reply, err := clients[idx].Set(ctx, &pb.SetRequest{Key: key})
			check(err)
			if reply.GetStatus() == statusCode.Moved {
				fmt.Printf("%s : %d -> %d ... first moved\n", key, idx, reply.GetTarget())
				fmt.Printf("%s : server next, level : %d %d \n", key, reply.GetNext(), reply.GetLevel())
				fmt.Printf("%s : local next,  level : %d %d \n", key, atomic.LoadInt64(&meta.next), atomic.LoadInt64(&meta.level))
				// 更新元数据
				// atomic.StoreInt64(&meta.next, reply.GetNext())
				// atomic.StoreInt64(&meta.level, reply.GetLevel())
				fmt.Printf("%s : local next,  level : %d %d ... after store\n", key, atomic.LoadInt64(&meta.next), atomic.LoadInt64(&meta.level))
				// 更新后元数据
				idx = hashFunc(key)
				if idx != reply.GetTarget() {
					panic(fmt.Sprintf("%s : server idx, local idx : %d, %d\n", key, reply.GetTarget(), idx))
				}
				reply, err = clients[idx].Set(ctx, &pb.SetRequest{Key: key})
				check(err)
				if reply.GetStatus() == statusCode.Moved {
					panic(fmt.Sprintf("%s : %d -> %d second moved, failed.\n", key, idx, reply.GetTarget()))
				}
			} else if reply.GetStatus() == statusCode.Failed {
				panic(fmt.Sprintf("%s : Err : %s", key, reply.GetErr()))
			}
		}
		// 分裂操作核心代码
		for index := 0; index < length; index++ {
			tempData := make([]string, 0) //学习redis来进行渐进性rehash
			db.slots[index].Lock()
			for _, key := range db.slots[index].data {
				if hashFunc(key) == secondIdx {
					set(secondIdx, key)
					count++
				} else {
					tempData = append(tempData, key)
				}
			}
			db.slots[index].data = tempData
			db.slots[index].Unlock()
		}

		// meta.Lock()
		// meta.serversStatus[*shardIdx] = serverStatus.Working
		// atomic.StoreInt32(&status, 1)
		// meta.Unlock()
		// atomic.StoreInt32(&status, 0)
		// fmt.Printf("%s : 分裂完成....\n", serversAddress[secondIdx])
		// 内部分裂完成，处理存起来的操作
		// fmt.Println("处理存起来的操作...")
		setCount := 0
		for _, operation := range operations {

			// 计算存起来的set操作有多少，看是否会出现积压过多导致二次分裂的问题
			switch operation[0] {
			case "set":
				setCount++
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
		fmt.Printf("%s 共处理 %d 个set积压操作\n", addresses[idx], setCount)

		// // 分裂完成，后置条件：将next指针往后移动
		// atomic.AddInt64(&meta.next, 1)
		// // 一轮分裂已经完成，将next指针归零
		// if atomic.LoadInt64(&meta.next) == int64(math.Pow(2, float64(atomic.LoadInt64(&meta.level))))*atomic.LoadInt64(&meta.hashSize) {
		// 	atomic.StoreInt64(&meta.next, 0)
		// }
		// fmt.Printf("Next : %d\n", atomic.LoadInt64(&meta.next))
		// // 该节点已经分裂过，所以针对该节点的数据，是需要将level值加1后重新计算的
		// atomic.AddInt64(&meta.level, 1)
	}

	// mutex.Unlock()
	splitInside()

	// mutex.Lock()
	// meta.version++
	atomic.AddInt64(&meta.version, 1)
	// 分裂过后，将该节点的阈值降为原来的4/5
	maxKeyOld := int64(float64(atomic.LoadInt64(&maxKey)) * 0.8)
	atomic.StoreInt64(&maxKey, maxKeyOld)
	// syncConf()
	// mutex.Unlock()
	// if full {
	// 	canSplit = false
	// 	syncConf()
	// 	return &pb.SplitReply{Status: statusCode.Ok, Result: 0, Full: true}, nil
	// }

	return &pb.SplitReply{Status: statusCode.Ok, Version: atomic.LoadInt64(&meta.version)}, nil
}

func (s *server) Scan(ctx context.Context, in *pb.ScanRequest) (*pb.ScanReply, error) {

	count := db.len
	meta.RLock()
	result := meta.serversAddress[*shardIdx] + "\n"
	meta.RUnlock()
	// result += fmt.Sprintf("%s", db.scan())
	result += db.scan()
	return &pb.ScanReply{Result: result, Status: statusCode.Ok, Count: int64(count), Version: atomic.LoadInt64(&meta.version)}, nil
}

func (s *server) SyncConf(ctx context.Context, in *pb.SyncConfRequest) (*pb.SyncConfReply, error) {
	meta.RLock()
	version := meta.version
	address := meta.serversAddress[*shardIdx]
	meta.RUnlock()
	if in.GetBegin() == *shardIdx {
		return &pb.SyncConfReply{Status: statusCode.Ok, Result: "同步结束", Version: version}, nil
	}
	fmt.Printf("%s : %s, 节点同步\n", time.Now().Format("2006-01-02 15:04:05"), address)
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
	}
	meta.Unlock()
	// 线性传递下去
	// 只要能将meta数据正确传递给下一个节点就行，如果id+1的节点不通就id+2...
	success := false
	target := *shardIdx + 1

	addresses := copyAddress(&meta)
	if target == int64(len(addresses)) {
		target = 0
	}
	for !success {
		serv := addresses[target]
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
		if target == int64(len(addresses)) {
			target = 0
		}
	}
	return &pb.SyncConfReply{Status: statusCode.Ok, Version: version}, nil
}

func (s *server) GetConf(ctx context.Context, in *pb.GetConfRequest) (*pb.GetConfReply, error) {

	addresses := copyAddress(&meta)
	var reply pb.GetConfReply
	reply.HashSize = atomic.LoadInt64(&meta.hashSize)
	reply.Level = atomic.LoadInt64(&meta.level)
	reply.Next = atomic.LoadInt64(&meta.next)
	reply.Status = statusCode.Ok
	for idx, addr := range addresses {
		var server pb.GetConfReply_ServConf
		server.Address = addr
		server.Idx = idx
		reply.Server = append(reply.Server, &server)
	}
	reply.Version = atomic.LoadInt64(&meta.version)
	reply.ShardIdx = atomic.LoadInt64(shardIdx)
	reply.ServerStatus = atomic.LoadInt32(&status)
	return &reply, nil
}

func (s *server) WakeUp(ctx context.Context, in *pb.WakeRequest) (*pb.WakeReply, error) {
	// 未知情况。。。
	if !atomic.CompareAndSwapInt32(&status, 1, 0) {
		log.Panicf("%d : 节点不在休眠状态\n", atomic.LoadInt64(shardIdx))
	}
	// 被唤醒并接受分裂的节点
	atomic.StoreInt64(&meta.next, in.GetNext())
	atomic.StoreInt64(&meta.level, in.GetLevel())
	version := atomic.LoadInt64(&meta.version)
	fmt.Printf("%s  节点被唤醒\n", time.Now().Format("2006-01-02 15:04:05"))
	return &pb.WakeReply{Status: statusCode.Ok, Version: version}, nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func initConf() {
	flag.Parse()

	statusCode.Init()
	errorCode.Init()
	serverStatus.Init()
	// check(err)
	operations = make([][]string, 0)
	meta.Lock()
	meta.hashSize = 4
	meta.serversAddress = make(map[int64]string)
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
		if idx == *shardIdx {
			atomic.StoreInt64(&maxKey, int64(ShardNodeConfs["max_key"].(float64)))
			if ShardNodeConfs["status"].(string) == serverStatus.Working {
				atomic.StoreInt32(&status, 0)
			} else {
				fmt.Printf("%d : %s\n", *shardIdx, serverStatus.Working)
				fmt.Printf("%d : %s\n", *shardIdx, ShardNodeConfs["status"].(string))
				atomic.StoreInt32(&status, 1)
			}
		}
		// fmt.Println(idx, ip, port)
	}
	meta.level = 0
	meta.Unlock()
}

func hashFunc(key string) int64 {
	posCRC16 := int64(crc16.Checksum([]byte(key), crc16.IBMTable))
	pos := posCRC16 % (int64(math.Pow(2.0, float64(atomic.LoadInt64(&meta.level))) * float64(atomic.LoadInt64(&meta.hashSize))))
	// if pos < atomic.LoadInt64(&meta.next) { // 分裂过了的
	// 	pos = posCRC16 % (int64(math.Pow(2.0, float64(atomic.LoadInt64(&meta.next)+1))) * atomic.LoadInt64(&meta.hashSize))
	// }
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
	addresses := copyAddress(&meta)
	fmt.Printf("db:\n%s\n", db.scan())
	// fmt.Printf("servers:\n%s\n", serversAddress)
	for idx, addr := range addresses {
		fmt.Printf("%d %s\t", idx, addr)
	}
	fmt.Println()
	// fmt.Printf("status:\n%s\n", serversStatus)
	fmt.Println()

	// fmt.Printf("maxKey:\n%s\n", serversMaxKey)
	// for idx, status := range meta0.serversStatus {
	// 	fmt.Printf("%d %s\t", idx, status)
	// }
	fmt.Println()
	fmt.Printf("shardIdx:\n%d\n", *shardIdx)
	fmt.Printf("next:\n%d\n", atomic.LoadInt64(&meta.next))
	fmt.Printf("level:\n%d\n", atomic.LoadInt64(&meta.level))
	fmt.Printf("hashSize:\n%d\n", atomic.LoadInt64(&meta.hashSize))
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

func copyAddress(m *Meta) map[int64]string {
	meta.RLock()
	addresses := make(map[int64]string, len(meta.serversAddress))
	for k, v := range meta.serversAddress {
		addresses[k] = v
	}
	meta.RUnlock()
	return addresses
}
