package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

func main() {
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
