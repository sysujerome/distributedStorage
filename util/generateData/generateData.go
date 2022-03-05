package main

import (
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	operators   = []string{"get", "set", "del"}
	// length      = flag.Int("length", 100000, "the number of cases")
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
	filename := filepath.Join(path, "data/test_data")
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	wirter, err := os.Create(filepath.Join(path, "data/data2"))
	check(err)
	defer wirter.Close()

	data, err := ioutil.ReadAll(f)
	check(err)
	kvs := strings.Split(string(data), "\n")
	newline := ""
	for i := 0; i+1 < len(kvs); i += 2 {
		wirter.WriteString(newline + kvs[i] + " " + kvs[i+1])
		// fmt.Println(newline + kvs[i] + " " + kvs[i+1])
		newline = "\n"
	}
}
