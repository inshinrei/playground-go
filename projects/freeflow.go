package projects

//Input: words = ["bella","label","roller"]
//Output: ["e","l","l"]
//
//Input: words = ["cool","lock","cook"]
//Output: ["c","o"]

func intersection(s1 map[string]bool, s2 map[string]bool) map[string]bool {
	var intersect map[string]bool
	if len(s1) > len(s2) {
		s2, s1 = s1, s2
	}
	for k, _ := range s1 {
		if s2[k] {
			intersect[k] = true
		}
	}
	return intersect
}

func commonChars(words []string) []string {
	var curr_set map[string]bool
	var common map[string]bool
	for _, char := range words[0] {
		common[char] = true
	}

	for word := range words {
		for _, char := range word {
			curr_set[char] = true
		}
		common = intersection(common, curr_set)
		curr_set = map[string]bool{}
	}
	var str string
	for char := range common {
		str += char
	}
	return str
}
