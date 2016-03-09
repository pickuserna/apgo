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

func testVariables() {
	var x, y int
	assertEqual(0, x)
	assertEqual(0, y)
}

func testForLoop() {
	result := 1
	var i int
	for {
		result *= 2
		i++
		if i >= 5 {
			break
		}
	}
	assertEqual(32, result)

	sum := 0
	for j := 0; j <= 5; j++ {
		sum += j
	}
	assertEqual(15, sum)
}

func main() {
	start := time.Now()
	testMath()
	testFunctions()
	testVariables()
	testForLoop()
	fmt.Println("Pass!")
	fmt.Println("Took ", time.Since(start))
}