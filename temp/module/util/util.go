package hello

import "fmt"

func SayHello(name string) {
	fmt.Println("Hello " + name)
}

type Obj struct {
	Name string
}
