package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	common "storage/util/common"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

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

func main() {
	confName := "conf_temp.json"
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
	count := 8
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
