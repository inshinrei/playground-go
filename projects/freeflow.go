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
