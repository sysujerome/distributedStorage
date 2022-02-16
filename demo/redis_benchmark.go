package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func Assert(check bool) {
	if !check {
		panic(check)
	}
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.1.128:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	currentPath, err := os.Getwd()
	check(err)
	filename := filepath.Join(currentPath, "data/benchmark_data")
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	sc := bufio.NewScanner(f)
	// data := make([][]string, 0)
	counter := 0
	start := time.Now()
	for sc.Scan() {
		operation := strings.Split(sc.Text(), " ")
		// fmt.Println(operation)
		opt := operation[0]
		switch opt {
		case "get":
			// continue
			_, _ = rdb.Get(ctx, operation[1]).Result()

		case "set":
			// continue
			res, _ := rdb.Set(ctx, operation[1], operation[2], 0).Result()
			Assert(res != " OK ")
			// fmt.Println("\"", res, "\"")
			// continue
		case "del":
			// continue
			_, err := rdb.Del(ctx, operation[1]).Result()
			check(err)
		}
		counter++
	}
	duration := time.Since(start)
	fmt.Printf("dealing with %d operations took %v Seconds\n", counter, duration.Seconds())
	// Output: key value
	// key2 does not exist
}
