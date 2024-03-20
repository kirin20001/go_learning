package data_structure_and_algorithm

func maxSubArray(nums []int) int {
	n := len(nums)
	dp := make([][2]int, n)
	dp[0][0] = nums[0]
	dp[0][1] = nums[0]

	for i := 1; i < n; i++ {
		dp[i][0] = max(dp[i-1][0], 0) + nums[i]
		dp[i][1] = max(dp[i-1][1], dp[i][0])
	}

	return dp[n-1][1]
}
