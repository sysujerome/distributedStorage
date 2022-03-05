package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// var input string
	// fmt.Scanf("%s", &input)
	// fmt.Printf("%s\n", input)
	// 输入一行带空格的字符串
	reader := bufio.NewReader(os.Stdin)
	input, _, _ := reader.ReadLine()
	fmt.Println(string(input))
}
