package main


import (
	crc16 "example.com/util/crc16"
	"fmt"
	"reflect"
)

func main()  {
	data := []byte("test")
	fmt.Println("begin")
	checksum := crc16.Checksum(data, crc16.IBMTable)

	fmt.Println("end")
	fmt.Println(reflect.TypeOf(checksum))
	fmt.Println(checksum)
}

