package data_structure_and_algorithm

/*
给你一个整数数组 prices ，其中 prices[i] 表示某支股票第 i 天的价格。
在每一天，你可以决定是否购买和/或出售股票。你在任何时候 最多 只能持有 一股 股票。你也可以先购买，然后在 同一天 出售。
返回 你能获得的 最大 利润 。
 */

func leetcode122_maxProfit(prices []int) int {
	n := len(prices)
	dp := make([][2]int, n)

	dp[0][0] = -prices[0]
	for i := 1; i < n; i++ {
		dp[i][0] = max(dp[i-1][1] - prices[i], dp[i-1][0])
		dp[i][1] = max(dp[i-1][0] + prices[i], dp[i-1][1])
	}
	return dp[n-1][1]
}

/*
动态规划
时间复杂度：O(n)
空间复杂度：O(n)
 */