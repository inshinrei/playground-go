package main

import (
	"fmt"
	"math/rand"
)

func main() {
	if n := rand.Intn(100); n == 0 {
		fmt.Println("random number is zero")
	} else if n > 50 {
		fmt.Println("random number is too big", n)
	} else {
		fmt.Println("random number is too big", n)
	}
}
