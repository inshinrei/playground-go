package freeflow

func main(words []string) {
	cnt := make([]int, 26)
	for i := range cnt {
		cnt[i] = 20000
	}
	for _, w := range words {
		t := make([]int, 26)
		for _, c := range w {
			t[c-'a']++
		}

	}
}
