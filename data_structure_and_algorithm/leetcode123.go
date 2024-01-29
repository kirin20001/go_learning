package data_structure_and_algorithm

import "math"

/*
给定一个数组，它的第 i 个元素是一支给定的股票在第 i 天的价格。
设计一个算法来计算你所能获取的最大利润。你最多可以完成 两笔 交易。
注意：你不能同时参与多笔交易（你必须在再次购买前出售掉之前的股票）。
 */
func leetcode123_maxProfit(prices []int) int {
	n := len(prices)
	dp1 := make([][2]int, n)
	dp2 := make([][2]int, n)
	dp1[0][0] = -prices[0]
	dp2[0][0] = math.MinInt8

	for i := 1; i < n; i++ {
		dp1[i][0] = max(dp1[i-1][0], -prices[i])
		dp1[i][1] = max(dp1[i-1][1], dp1[i][0] + prices[i])
		dp2[i][0] = max(dp2[i-1][0], dp1[i][1] - prices[i])
		dp2[i][1] = max(dp2[i-1][1], dp2[i][0] + prices[i])
	}

	return dp2[n-1][1]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

/*
动态规划
时间复杂度：O(n)
空间复杂度：O(n)
 */