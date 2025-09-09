package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// 解法一：使用map
func singleNumber1(nums []int) int {
	if len(nums) == 1 {
		return nums[0]
	}
	count := make(map[int]int)
	for i := 0; i < len(nums); i++ {
		count[nums[i]]++
	}
	for k, v := range count {
		if v == 1 {
			return k
		}
	}
	return 0
}

// 解法二：使用异或特性来找到出现单次的值
func singleNumber2(nums []int) int {
	single := 0
	for _, num := range nums {
		single ^= num
	}
	return single
}

// 解法一：转化成字符串使用双指针进行判断是否一致
func isPalindrome1(x int) bool {
	if x < 0 {
		return false
	}
	xStr := strconv.Itoa(x)
	for l, r := 0, len(xStr)-1; l < r; {
		if xStr[l] != xStr[r] {
			return false
		}
		l++
		r--
	}
	return true
}

// 解法二：不断将x值移到一个新变量里，直到x相等或者小于新变量，进行判断x是否与新变量相等
func isPalindrome2(x int) bool {
	if x < 0 || (x%10 == 0 && x != 0) {
		return false
	}

	revertedNumber := 0
	for x > revertedNumber {
		revertedNumber = revertedNumber*10 + x%10
		x /= 10
	}
	return x == revertedNumber || x == revertedNumber/10
}

// 解法一：使用栈，如果最终栈内有值则代表存在问题
func isValid1(s string) bool {
	stack := make([]string, len(s))
	countIndex := -1
	sNew := strings.Split(s, "")
	bracketsDict := map[string]string{
		"(": ")",
		"[": "]",
		"{": "}",
	}
	for i := 0; i < len(sNew); i++ {
		if countIndex == -1 || bracketsDict[stack[countIndex]] != sNew[i] {
			countIndex++
			stack[countIndex] = sNew[i]
		} else {
			stack[countIndex] = ""
			countIndex--
		}
	}
	return countIndex == -1
}

// 加入switch 进行便捷添加
func isValid2(s string) bool {
	stack := make([]string, len(s))
	countIndex := -1
	sNew := strings.Split(s, "")
	bracketsDict := map[string]string{
		"(": ")",
		"[": "]",
		"{": "}",
	}
	for i := 0; i < len(sNew); i++ {
		switch sNew[i] {
		case "(", "[", "{":
			countIndex++
			stack[countIndex] = sNew[i]
		case ")", "]", "}":
			if countIndex == -1 || bracketsDict[stack[countIndex]] != sNew[i] {
				countIndex++
				stack[countIndex] = sNew[i]
			} else {
				stack[countIndex] = ""
				countIndex--
			}
		}
	}
	return countIndex == -1
}

// 思路：从底位开始累加，加成功后终止。如果都等于9，则进位
func plusOne(digits []int) []int {
	n := len(digits)
	for i := n - 1; i >= 0; i-- {
		if digits[i] < 9 {
			digits[i]++
			return digits
		}
		digits[i] = 0
	}

	return append([]int{1}, digits...)
}

// 解法一：遍历整个数组，通过切片来进行追加数组间接删除该值
func removeDuplicates1(nums []int) int {
	for i := 1; i < len(nums); i++ {
		if nums[i] == nums[i-1] {
			nums = append(nums[:i], nums[i+1:]...)
			i--
		}
	}
	return len(nums)
}

// 解法二：双指针 只要和前面值不同则将该值复制到前面，之后进行切片
func removeDuplicates2(nums []int) int {
	if len(nums) <= 1 {
		return 1
	}

	low := 0
	for fast := 1; fast < len(nums); fast++ {
		if nums[fast] != nums[low] {
			nums[low+1] = nums[fast]
			low++
		}
	}
	nums = nums[:low+1]
	return len(nums)
}

// 按照每个区间的start进行排序，之后那当前区间的end 和后面区间的start 不断比较，如果发现start 大于end，则进行改变end。否则存储为新的区间
func merge(intervals [][]int) (ans [][]int) {
	slices.SortFunc(intervals, func(p, q []int) int { return p[0] - q[0] }) // 按照左端点从小到大排序
	for _, p := range intervals {
		m := len(ans)
		if m > 0 && p[0] <= ans[m-1][1] { // 可以合并
			ans[m-1][1] = max(ans[m-1][1], p[1]) // 更新右端点最大值
		} else { // 不相交，无法合并
			ans = append(ans, p) // 新的合并区间
		}
	}
	return
}

// 解法一：爆破
func twoSum1(nums []int, target int) []int {
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				return []int{i, j}
			}
		}
	}
	return []int{}
}

// 解法二：哈希解法
func twoSum2(nums []int, target int) []int {
	diffMap := make(map[int]int)
	for idx, num := range nums {
		if diffIdx, ok := diffMap[target-num]; ok {
			return []int{idx, diffIdx}
		}
		diffMap[num] = idx
	}
	return []int{}
}

func main() {
	// 只出现一次的数字
	fmt.Println(singleNumber1([]int{1, 1, 4, 3, 4}))
	fmt.Println(singleNumber2([]int{1, 1, 4, 3, 4}))
	// 回文数
	fmt.Println(isPalindrome1(121))
	fmt.Println(isPalindrome2(12321))
	// 有效括号
	fmt.Println(isValid1("()"))
	fmt.Println(isValid2("(])"))
	// 加一
	fmt.Println(plusOne([]int{1, 2, 3}))
	// 删除有序数组中的重复项
	fmt.Println(removeDuplicates1([]int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}))
	fmt.Println(removeDuplicates2([]int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}))
	// 合并区间
	fmt.Println(merge([][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}, {1, 6}, {8, 10}, {15, 18}, {4, 18}}))
	// 两数之和
	fmt.Println(twoSum1([]int{1, 3, 3}, 6))
	fmt.Println(twoSum2([]int{1, 3, 3}, 6))
}
