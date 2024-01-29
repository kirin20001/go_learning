package data_structure_and_algorithm

/*
给你一个 非严格递增排列 的数组 nums ，请你 原地 删除重复出现的元素，使每个元素 只出现一次 ，返回删除后数组的新长度。
元素的 相对顺序 应该保持 一致 。然后返回 nums 中唯一元素的个数。
考虑 nums 的唯一元素的数量为 k ，你需要做以下事情确保你的题解可以被通过：
更改数组 nums ，使 nums 的前 k 个元素包含唯一元素，并按照它们最初在 nums 中出现的顺序排列。
nums 的其余元素与 nums 的大小不重要。
返回 k 。
 */
func leetcode26_removeDuplicates(nums []int) int {
	removeFlag:= 1
	for i := 0; i < len(nums)-1; i++ {
		if nums[i] != nums[i+1] {
			nums[removeFlag] = nums[i+1]
			removeFlag++
		}
	}
	return removeFlag
}

/*
使用双指针，指针1用于生成最终输出的数组；指针2类似一个长度为2的滑动窗口，检查原数组是否有新元素出现。
时间复杂度：O(n)
空间复杂度：O(1)
 */