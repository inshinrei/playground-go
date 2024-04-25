package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	c, _ := net.Dial("tcp", "127.0.0.1:8080")
	_, err := fmt.Fprintf(c, "GET / HTTP/1.0\r\n\r\n")
	if err != nil {
		return
	}
	s, _ := bufio.NewReader(c).ReadString('\n')
	fmt.Println(s)
}
