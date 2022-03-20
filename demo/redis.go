package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.1.128:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set(ctx, "jerome", "cliffia", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "jerome").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("jerome", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}

	fmt.Println("删除键jerome...")
	val3, err := rdb.Del(ctx, "jerome").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(val3)

	fmt.Println("删除键jerome之后")
	val4, err := rdb.Get(ctx, "jerome").Result()
	if err == redis.Nil {
		fmt.Println("jerome does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("jerome", val4)
	}

	err = rdb.Set(ctx, "jerome", "cliffia", 0).Err()
	if err != nil {
		panic(err)
	}
	err = rdb.Set(ctx, "a", "cliffia", 0).Err()
	if err != nil {
		panic(err)
	}
	err = rdb.Set(ctx, "aa", "cliffia", 0).Err()
	if err != nil {
		panic(err)
	}
	// Output: key value
	// key2 does not exist
}
