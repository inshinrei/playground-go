package main

import "sort"

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
