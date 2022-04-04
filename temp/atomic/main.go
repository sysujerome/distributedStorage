package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var target int64 = 0

func main() {
	timing()
}

// func main() {
// 	type Meta struct {
// 		serversAddress map[int64]string
// 		serversStatus  map[int64]string
// 		serversMaxKey  map[int64]int64
// 		next           int64
// 		level          int64
// 		hashSize       int64
// 		version        int64
// 		sync.RWMutex
// 	}
// 	worker1 := func() {
// 		var m Meta
// 		m.serversMaxKey = make(map[int64]int64, 0)
// 		time.Sleep(time.Second)
// 		m.serversMaxKey[1] = 2
// 		fmt.Printf("%d", m.serversMaxKey[1])
// 		// time.Sleep(time.Second)
// 	}

// 	// var wg sync.WaitGroup
// 	// wg.Add(2)
// 	// go worker(&wg)
// 	// go worker(&wg)
// 	// wg.Wait()
// 	for i := 0; i < 1000; i++ {
// 		go worker1()

// 	}
// 	// time.Sleep(time.Millisecond * 100)
// 	go worker1()
// 	// time.Sleep(time.Millisecond * 100)
// 	go worker1()
// 	// time.Sleep(time.Millisecond * 100)
// 	go worker1()
// 	// time.Sleep(time.Millisecond * 100)
// 	go worker1()
// 	// time.Sleep(time.Millisecond * 100)
// 	go worker1()
// 	// time.Sleep(time.Millisecond * 100)
// 	go worker1()
// 	// time.Sleep(time.Millisecond * 100)
// 	go worker1()
// 	// time.Sleep(time.Millisecond * 100)
// 	go worker1()
// 	// time.Sleep(time.Millisecond * 100)
// 	go worker1()
// 	// time.Sleep(time.Millisecond * 100)
// 	go worker1()
// 	// time.Sleep(time.Millisecond * 100)
// 	go worker1()
// 	time.Sleep(time.Second * 2)
// }
func worker(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i <= 10; i++ {
		atomic.AddInt64(&target, int64(i))
	}
	fmt.Printf("%d\n", target)
}

func slice() {
	a := make([]int, 5)
	a = append(a, 1)
	a = append(a, 2)
	a = append(a, 3)
	a = append(a, 4)
	a = append(a, 5)
	fmt.Println(a, len(a), cap(a))
}
