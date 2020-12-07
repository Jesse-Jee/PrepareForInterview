package main

import "fmt"

func insertSort(nums []int) {

	for i := 1; i < len(nums); i++ {
		pre, cur := i-1, nums[i]

		for pre >= 0 && nums[pre] > cur {
			nums[pre+1] = nums[pre]
			pre--
		}
		nums[pre+1] = cur
	}
}

func main() {
	nums := []int{9, 8, 6, 7, 4, 5, 1, 3, 2}
	insertSort(nums)
	fmt.Println(nums)
}
