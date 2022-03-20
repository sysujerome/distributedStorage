package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	ip   = flag.String("ip", "192.168.1.128", "the ip of the shard")
	port = flag.String("port", "50050", "the port of the shard")
)

func main() {
	fmt.Printf("%s : %s\n", "ip", *ip)
	fmt.Printf("%s : %s\n", "port", *port)
	// flag.Parse()
	initConf()
	fmt.Printf("%s : %s\n", "ip", *ip)
	fmt.Printf("%s : %s\n", "port", *port)
}

func initConf() {
	flag.Parse()
	*port = "22222"
	fmt.Printf("%s : %s\n", "port", *port)
	flag.Parse()
	fmt.Printf("%s : %s\n", "port", *port)
	fmt.Printf("%s\n", os.Args[1])
}

/*
using :
	$ go build -o flag ./test/flag.go
	$ ./flag --port 100
*/
