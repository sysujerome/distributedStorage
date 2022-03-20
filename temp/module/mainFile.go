package main

import (
	"fmt"
	"hello"
)

func main() {
	a := new(hello.Obj)
	fmt.Println("jerome")
	a.Name = "cliffia"
	hello.SayHello(a.Name)
}
