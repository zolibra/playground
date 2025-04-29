package main

import (
	"fmt"
)

func main() {
	a := 10
	b := 0

	result := divide(a, b)

	fmt.Println("Result of division is:", result)
}

func divide(x int, y int) int {
	return x / y
}
