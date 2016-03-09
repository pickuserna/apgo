package main

import "fmt"

func addOne(x int) int {
	return x + 1
}

func main() {
	x := 1 + 1
	x = addOne(x)
	if x > 2 {
		fmt.Println("x is greater than 2")
	} else {
		fmt.Println("x is not greater than 2")
	}
}