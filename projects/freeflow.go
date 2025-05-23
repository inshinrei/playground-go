package main

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

func main() {
	nums := []int{5, 6, 3, 4, 1}
	maximumDifference(nums)
}

func maximumDifference(nums []int) int {
	var ans = -1
	var mn = nums[0]

	for _, v := range nums {
		fmt.Println(v, mn)
		if v > mn {
			ans = max(ans, v-mn)
		}
		mn = min(mn, v)
	}
	return ans
}

func maximumImportance(n int, roads [][]int) (ans int64) {
	deg := make([]int, n)
	for _, r := range roads {
		deg[r[0]] += 1
		deg[r[1]] += 1
	}
	sort.Ints(deg)
	for i := 0; i < n; i++ {
		ans += int64((i + 1) * deg[i])
	}
	return ans
}

func checkOverlap(radius, xCenter, yCenter, x1, y1, x2, y2 int) bool {
	distanceX := distanceToEdge(x1, x2, xCenter)
	distanceY := distanceToEdge(y1, y2, yCenter)
	return distanceX*distanceX+distanceY*distanceY <= radius*radius
}

func distanceToEdge(minEdge, maxEdge, center int) int {
	if minEdge <= center && center <= maxEdge {
		return 0
	}
	if center < minEdge {
		return minEdge - center
	}
	return center - maxEdge
}

type Person struct {
	name string
	age  uint8
}

func changePerson(person **Person) {
	*person = &Person{
		name: "A",
		age:  25,
	}
}

func per() {
	person := &Person{
		name: "B",
		age:  2,
	}
	changePerson(&person)
}

// task 1
func processData(v int) int {
	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	return v * 2
}

func processParallel(in, out chan int, numWorkers int) {}

func runTask1() {
	in := make(chan int)
	out := make(chan int)

	go func() {
		for i := range 10 {
			in <- i
		}
		close(in)
	}()

	start := time.Now()
	processParallel(in, out, 5)
	for v := range out {
		fmt.Println("v =", v)
	}
	fmt.Println("main duration:", time.Since(start))
}

// chan
func writer() <-chan int {
	ch := make(chan int)
	wg := &sync.WaitGroup{}

	wg.Add(3)
	go func() {
		defer wg.Done()
		for i := range 10 {
			ch <- i
		}
	}()

	go func() {
		defer wg.Done()
		for i := range 10 {
			ch <- i + 5
		}
	}()

	go func() {
		defer wg.Done()
		for i := range 10 {
			ch <- i + 15
		}
	}()

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}

func ch() {
	ch1 := writer()

	for d := range ch1 {
		fmt.Println("ch1 =", d)
	}
}

func doubler(x int) int {
	time.Sleep(5 * time.Second)
	return x * 2
}

//func writer1() <-chan int {
//	ch := make(chan int)

//ctx, cancel = context.WithTimeout(context.Background, 5*time.Second)
//}

func cont1() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ch := make(chan int)

	go func() {
		for i := range 1000 {
			select {
			case ch <- i:
			case <-ctx.Done():
				break
			}
		}
	}()

	for {
		select {
		case v, ok := <-ch:
			if !ok {
				return
			}
			fmt.Println("ch1 =", v)
		case <-ctx.Done():
			return
		}
	}
}

// task 2
func randomTimework() {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Second)
}

func predictableTimework() {
	ch := make(chan struct{})

	go func() {
		randomTimework()
		close(ch)
	}()

	select {
	case <-ch:
	case <-time.After(3 * time.Second):
		panic("timed out")
	}
}
