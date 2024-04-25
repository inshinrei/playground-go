package main

import (
	"fmt"
	"net/http"
)

func hello(res http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprint(res, "hey")
	if err != nil {
		return
	}
}

func main() {
	http.HandleFunc("/", hello)
	err := http.ListenAndServe("localhost:4040", nil)
	if err != nil {
		return
	}
}
