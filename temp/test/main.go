package main

import "fmt"

func main() {
	var maxKey int64 = 12345
	fmt.Printf("before : %d\n", maxKey)

	maxKey = int64(float64(maxKey) * 1.25)

	fmt.Printf("after : %d\n", maxKey)
	maxKey = int64(float64(maxKey) * 0.8)

	fmt.Printf("reset : %d\n", maxKey)

	nums := []float32{
		139.20,
		138.00,
		84.35,
		89.20,
		350.06,
	}
	var sum float32 = 0
	for _, n := range nums {
		sum += n
	}
	fmt.Println(sum)
}
