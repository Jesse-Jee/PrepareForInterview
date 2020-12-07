package main

import "fmt"

func quickSort(nums []int, begin, end int) {
	if begin > end {
		return
	}

	pivot := partition(nums, begin, end)

	quickSort(nums, begin, pivot-1)
	quickSort(nums, pivot+1, end)
}

func partition(nums []int, begin, end int) int {
	pivot := end
	i := begin

	for j := begin; j < end; j++ {
		if nums[j] < nums[pivot] {
			nums[i], nums[j] = nums[j], nums[i]
			i++
		}
	}

	nums[i], nums[pivot] = nums[pivot], nums[i]
	return i
}

func main() {
	nums := []int{3, 9, 8, 5, 6, 4, 2, 7, 1}
	quickSort(nums, 0, len(nums)-1)
	fmt.Println(nums)
}
