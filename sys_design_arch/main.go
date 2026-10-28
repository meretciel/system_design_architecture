package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, world!")

	x := []int{1, 2, 3, 4}
	num := copy(x[1:], x[:3])
	fmt.Println(x, num)

	for i := 0; i < 10; i++ {
		fmt.Println(i)
	}
}
