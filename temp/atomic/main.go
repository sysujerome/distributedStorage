package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var target int64 = 0

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go worker(&wg)
	go worker(&wg)
	wg.Wait()
}
func worker(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i <= 10; i++ {
		atomic.AddInt64(&target, int64(i))
	}
	fmt.Printf("%d\n", target)
}
