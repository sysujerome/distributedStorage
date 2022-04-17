package main

import "os"

func main() {
	filename := "test_file"
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	file.WriteString("sdjakjsd\n")
	file.Close()

}
