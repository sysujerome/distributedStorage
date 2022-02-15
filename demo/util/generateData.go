package main

import (
	"math/rand"
	"os"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	operators   = []string{"get", "set", "del"}
)

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// you need to ensure that the base directory exists, because os.Create does not handle it.
	// err := os.WriteFile("pangjinglongtmpfile", []byte(RandStringRunes(10)), 0644)
	f, err := os.Create("benchmark_data")
	newline := ""
	var key, value string
	check(err)
	defer f.Close()
	for i := 0; i < 1000; i++ {
		opt := operators[rand.Intn(len(operators))]
		switch opt {
		case "get":
			key = RandStringRunes(10)
			f.WriteString(newline + opt + " " + key)
		case "set":
			key = RandStringRunes(10)
			value = RandStringRunes(10)
			f.WriteString(newline + opt + " " + key + " " + value)
		case "del":
			key = RandStringRunes(10)
			f.WriteString(newline + opt + " " + key)
		}
		newline = "\n"
	}
}
