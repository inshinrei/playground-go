package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

const (
	unlocked = false
	locked   = true
)

type BrokenMutex struct {
	state atomic.Bool
}

func (m *BrokenMutex) Lock() {
	for !m.state.CompareAndSwap(unlocked, locked) {
	}
}

func (m *BrokenMutex) Unlock() {
	m.state.Store(unlocked)
}

const n = 1000

func main() {
	var bm BrokenMutex
	wg := &sync.WaitGroup{}
	wg.Add(n)

	value := 0
	for range n {
		go func() {
			defer wg.Done()
			bm.Lock()
			value++
			bm.Unlock()
		}()
	}

	wg.Wait()
	fmt.Printf("%d\n", value)
}
