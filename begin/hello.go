package main

import (
	"fmt"
	"time"
)

func count() {
	for i := 0; i < 5; i++ {
		fmt.Println(i)
		time.Sleep(1 * time.Millisecond)
	}
}

func main() {
	go count()
	time.Sleep(time.Millisecond * 2)
	fmt.Println("oops")
	time.Sleep(time.Millisecond * 5)
}
