package main

import (
	"fmt"
	"time"
)

func fib(n int) int {
	if n < 2 {
		return 1
	}
	return fib(n - 1) + fib(n - 2)
}

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
	n := 26
	start := time.Now()
	fmt.Println("fib(", n, ") is ", fib(n))
	fmt.Println("Took ", time.Since(start))
}