package main

import "fmt"

func getName() string {
	return "someone"
}

func main() {
	name := getName()
	fmt.Println(name)
}
