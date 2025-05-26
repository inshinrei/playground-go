package main

import (
	"context"
	"fmt"
	"runtime"
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
		runtime.Gosched()
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

func fanIn(ctx context.Context, chans []chan int) chan int {
	outCh := make(chan int)

	go func() {
		wg := &sync.WaitGroup{}
		wg.Add(len(chans))
		for _, ch := range chans {
			go func() {
				defer wg.Done()
				for {
					select {
					case v, ok := <-ch:
						if !ok {
							return
						}
						select {
						case outCh <- v:
						case <-ctx.Done():
							return
						}
					case <-ctx.Done():
						return
					}
				}
			}()
		}
		wg.Wait()
		close(outCh)
	}()

	return outCh
}

func fanOut(ch chan int, n int) []chan int {
	outChs := make([]chan int, n)
	for i := range n {
		outChs[i] = make(chan int)
	}
	return outChs
}
