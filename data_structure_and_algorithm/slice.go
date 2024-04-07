package data_structure_and_algorithm

import (
	"fmt"
	"math"
	"sort"
)

func slictest() {
	m := make(map[int] int)
	m[1] = 1
	m[2] = 2
	m[3] = 3
	fmt.Println(len(m))

	var a int
	a = math.MaxInt - 1  + 1
	fmt.Println(a)

	b := []int{0, 1, 2, 3}
	fmt.Println(b[1:4])

	s := "abcde"
	fmt.Println(s[1:5])
}

type Node struct {
	val int
}

func sliceSort() {
	nums := []int{2,3,1,5,4,9,6,0,7}
	nums2 := []string{"a", "p", "b"}
	nums3 := [][]int{{9,9}, {1,1}, {1,5}, {3,5}}
	sort.Slice(nums, func(i, j int) bool {return nums[i] < nums[j]})
	fmt.Println(nums)
	sort.Slice(nums2, func(i, j int) bool {return nums2[i] < nums2[j]})
	fmt.Println(nums2)
	sort.Slice(nums3, func(i, j int) bool {return nums3[i][0] < nums3[j][0]})
	fmt.Println(nums3)

	fmt.Println(nums)
}

func copySlice(n []int) []int {
	return n
}
