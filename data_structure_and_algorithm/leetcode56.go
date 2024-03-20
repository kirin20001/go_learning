package data_structure_and_algorithm

import "sort"

func merge(intervals [][]int) [][]int {
	nums := []int{}
	for _, v := range intervals {
		nums = append(nums, v[0]*10)
		nums = append(nums, v[1]*10+1)
	}
	sort.Ints(nums)
	var flag int
	ans := make([][]int, 0)
	temp := []int{0, 0}
	for _, v := range nums {
		if v%10 == 0{
			if flag == 0 {
				temp[0] = v/10
			}
			flag++
		}

		if v%10 == 1 {
			flag--
			if flag == 0 {
				temp[1] = v/10
				ans = append(ans, temp)
				temp = []int{0, 0}
			}
		}
	}
	return ans
}