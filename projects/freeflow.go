package main

import (
	"fmt"
	"math/rand"
)

func main() {
	s := make([]int, 100)
	for i := 0; i < 100; i++ {
		s[i] = rand.Intn(100)
	}
	fmt.Println(s)

	for _, v := range s {
		switch {
		case v%2 == 0 && v%3 == 0:
			fmt.Println("six")
		case v%2 == 0:
			fmt.Println("two")
		case v%3 == 0:
			fmt.Println("three")
		default:
			fmt.Println("un")
		}
	}
}
