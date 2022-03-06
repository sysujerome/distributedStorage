package main

import (
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
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
	generate_test_data()
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
