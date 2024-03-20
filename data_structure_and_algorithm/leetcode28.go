package data_structure_and_algorithm

import "fmt"

func strStr(haystack string, needle string) int {
	n := len(needle)
	if n == 0 {
		return 0
	}
	next := make([]int, len(needle))
	getNext(next, needle)
	fmt.Println("strStr", next)
	j := 0
	for i := 0; i < len(haystack); i++ {
		for j > 0 && haystack[i] != needle[j] {
			j = next[j]
		}
		if haystack[i] == needle[j] {
			j++
		}
		if j == n {
			return i + 1 - n
		}
	}
	return -1
}

func getNext(next []int, s string) {
	j := 0
	next[0] = 0
	for i := 1; i < len(s); i++ {
		fmt.Println(i, next)
		for j > 0 && s[i] != s[j] {
			j = next[j-1]
		}
		if s[i] == s[j] {
			j++
		}
		next[i] = j
	}
}
