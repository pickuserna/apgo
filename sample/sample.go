package main

import (
	"fmt"
	"time"
)

func assertEqual(a interface{}, b interface{}) {
	if a != b {
		panic(fmt.Sprint("Expected ", a, ", but got ", b))
	}
}

func fib(n int) int {
	if n < 2 {
		return 1
	}
	return fib(n - 1) + fib(n - 2)
}

func addOne(x int) int {
	return x + 1
}

func testMath() {
	assertEqual(2, 1 + 1)
}

func testFunctions() {
	assertEqual(5, fib(4))
	assertEqual(2, addOne(1))
}

func main() {
	start := time.Now()
	testMath()
	testFunctions()
	fmt.Print("Pass!")
	fmt.Println("Took ", time.Since(start))
}