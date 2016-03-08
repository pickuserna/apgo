package main

import "fmt"

func main() {
	x := 1 + 1
	x += 1
	if x > 2 {
		fmt.Println("x is greater than 2")
	} else {
		fmt.Println("x is not greater than 2")
	}
}