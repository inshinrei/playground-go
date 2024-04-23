package main

import "testing"

func TestName(t *testing.T) {
	name := getName()
	if name != "someone" {
		t.Error("response from getName is unexpected")
	}
}
