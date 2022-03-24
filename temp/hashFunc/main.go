package main

import (
	"fmt"
	"math"
	crc16 "storage/util/crc16"
	"time"
)

func hashFunc(key string) int64 {
	level := 0
	var next int64 = 0
	var hashSize int64 = 4
	posCRC16 := int64(crc16.Checksum([]byte(key), crc16.IBMTable))
	pos := posCRC16 % (int64(math.Pow(2.0, float64(level))) * hashSize)
	if pos < next { // 分裂过了的
		pos = posCRC16 % (int64(math.Pow(2.0, float64(level+1))) * hashSize)
	}
	return pos
}
func main() {
	key := "sdaldjkas"
	fmt.Scanf("%s", &key)
	for i := 0; i < 10; i++ {
		fmt.Println(hashFunc(key))
		fmt.Println(time.Now().Nanosecond())
	}
}
