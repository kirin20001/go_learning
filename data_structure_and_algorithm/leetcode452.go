package data_structure_and_algorithm

import "sort"

func findMinArrowShots(points [][]int) int {
	sort.Slice(points, func(i, j int) bool {return points[i][1] < points[j][1]})
	count := 0
	i := 0
	for i < len(points) {
		right := points[i][1]
		i++
		for i < len(points) && points[i][0] <= right {
			i++
		}
		count++
	}
	return count
}
