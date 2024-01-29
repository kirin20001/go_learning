package data_structure_and_algorithm

/*
88. 合并两个有序数组
给你两个按 非递减顺序 排列的整数数组 nums1 和 nums2，另有两个整数 m 和 n ，分别表示 nums1 和 nums2 中的元素数目。
请你 合并 nums2 到 nums1 中，使合并后的数组同样按 非递减顺序 排列。
注意：最终，合并后数组不应由函数返回，而是存储在数组 nums1 中。
为了应对这种情况，nums1 的初始长度为 m + n，其中前 m 个元素表示应合并的元素，后 n 个元素为 0 ，应忽略。nums2 的长度为 n 。
*/

func leetcode88_merge(nums1 []int, m int, nums2 []int, n int)  {
	i, j, flag := m-1, n-1, len(nums1)-1
	for i >= 0 && j >= 0 {
		if nums2[j] >= nums1[i] {
			nums1[flag] = nums2[j]
			j--
			flag--
			continue
		}
		if nums1[i] > nums2[j] {
			nums1[flag] = nums1[i]
			i--
			flag--
			continue
		}
	}

	if j >= 0 {
		for a := 0; a <=j; a ++ {
			nums1[a] = nums2[a]
		}
	}
}

/*
合并两个有序数组，从后往前遍历保证遍历一边后能够顺序写入nums1。
时间复杂度: O(n)
空间复杂度: O(1)
 */