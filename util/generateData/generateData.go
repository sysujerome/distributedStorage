package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	common "storage/util/common"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	// operators   = []string{"get", "set", "del"}
	// length      = flag.Int("length", 100000, "the number of cases")
)

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b) + strconv.Itoa(time.Now().Nanosecond())
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// func main() {
// 	// you need to ensure that the base directory exists, because os.Create does not handle it.
// 	// err := os.WriteFile("pangjinglongtmpfile", []byte(RandStringRunes(10)), 0644)
// 	flag.Parse()
// 	path, _ := os.Getwd()
// 	filename := filepath.Join(path, "/data/benchmark_data")
// 	f, err := os.Create(filename)
// 	newline := ""
// 	var key, value string
// 	check(err)
// 	defer f.Close()
// 	length, err := strconv.ParseInt(os.Args[1], 0, 64)
// 	check(err)
// 	for i := 0; i < int(length); i++ {
// 		opt := operators[rand.Intn(len(operators))]
// 		switch opt {
// 		case "get":
// 			key = RandStringRunes(10)
// 			f.WriteString(newline + opt + " " + key)
// 		case "set":
// 			key = RandStringRunes(10)
// 			value = RandStringRunes(10)
// 			f.WriteString(newline + opt + " " + key + " " + value)
// 		case "del":
// 			key = RandStringRunes(10)
// 			f.WriteString(newline + opt + " " + key)
// 		}
// 		newline = "\n"
// 	}
// }

func main() {
	// generate_test_data()
	if len(os.Args) < 2 {
		panic("Too less argument...")
	}
	count, err := strconv.Atoi(os.Args[1])
	check(err)
	generate_conf(count)
}

func generate_test_data() {
	path, _ := os.Getwd()
	wirter, err := os.Create(filepath.Join(path, "data/data2"))
	check(err)
	defer wirter.Close()

	length := 10 * 10000
	newline := ""
	for i := 0; i < length; i++ {
		key := RandStringRunes(10)
		value := RandStringRunes(10)
		wirter.WriteString(newline + key + " " + value)
		newline = "\n"
	}
}

func generate_conf(count int) {
	// {
	// 	"next":0,
	// 	"level":0,
	// 	"shard_confs":[
	// 		{
	// 			"shard_idx" :0,
	// 			"shard_node_confs" : {
	// 				"shard_node_0":  {
	// 					"ip":"192.168.1.128",
	// 					"base_port":50050,
	// 					"status":"working",
	// 					"max_key":1000
	// 				}
	// 			}
	// 		},
	type ShardNodeConf struct {
		IP       string `json:"ip"`
		BasePort int    `json:"base_port"`
		Status   string `json:"status"`
		MaxKey   int64  `json:"max_key"`
	}
	type shardConfs struct {
		ShardIdx        int            `json:"shard_idx"`
		ShardNodeName   string         `json:"shard_node_name"`
		Shard_node_conf *ShardNodeConf `json:"shard_node_confs"`
	}
	type conf struct {
		Next        int           `json:"next"`
		Level       int           `json:"level"`
		Shard_confs []*shardConfs `json:"shard_confs"`
	}

	confName := "conf.json"
	path, err := os.Getwd()
	check(err)
	fileName := filepath.Join(path, confName)
	writer, err := os.Create(fileName)
	check(err)
	defer writer.Close()

	var statusCode common.Status
	var errorCode common.Error
	var serverStatus common.ServerStatus
	statusCode.Init()
	errorCode.Init()
	serverStatus.Init()
	shards := make([]*shardConfs, 0)
	for i := 0; i < count; i++ {
		shard := new(shardConfs)
		shard.ShardIdx = i
		shard.ShardNodeName = fmt.Sprintf("shard_node_%d", i)
		shardConf := new(ShardNodeConf)
		shardConf.IP = "192.168.1.128"
		shardConf.BasePort = 50050 + i
		shardConf.MaxKey = 1000
		// Working  string
		// Sleep    string
		// Spliting string
		// Full     string
		shardConf.Status = serverStatus.Sleep
		if i < 4 {
			shardConf.Status = serverStatus.Working
		}
		shard.Shard_node_conf = shardConf
		shards = append(shards, shard)
	}

	confData := conf{
		Next:        0,
		Level:       0,
		Shard_confs: shards,
	}
	confJson, err := json.Marshal(confData)
	var confStr bytes.Buffer
	_ = json.Indent(&confStr, []byte(confJson), "", "	")
	check(err)
	writer.Write(confStr.Bytes())
	// return

}
