package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var (
	shard_idx       int
	shard_node_name string = "shard_node_0"
)

func main() {
	var configure map[string]interface{}
	fileData, err := ioutil.ReadFile("tmp")
	if err != nil {
		panic(err)
	}
	// fmt.Printf("%s", data)
	json.Unmarshal(fileData, &configure)
	confs := configure["shard_confs"].([]interface{}) // shard_confs
	for _, v := range confs {
		value := v.(map[string]interface{})
		shard_idx = int(value["shard_idx"].(float64))
		shard_node_confs := value["shard_node_confs"].(map[string]interface{})
		shard_content, found := shard_node_confs[shard_node_name].(map[string]interface{})
		if !found {
			return
		}
		for k, v := range shard_content {
			fmt.Printf("%s, %s\n", k, v)
		}
	}
}
