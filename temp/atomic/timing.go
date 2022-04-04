package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

const num int32 = 1000000

func timing() {
	var target int32 = -1
	var i int32 = 0
	fmt.Println("store 100w operation:")
	start := time.Now()
	for ; i < num; i++ {
		atomic.StoreInt32(&target, i)
	}
	elapse := time.Since(start)
	fmt.Println(elapse)

	fmt.Printf("\nswap 100w operation: \n")
	start = time.Now()
	i = 0
	for ; i < num; i++ {
		atomic.SwapInt32(&target, i)
	}
	elapse = time.Since(start)
	fmt.Println(elapse)
}
