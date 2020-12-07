package main

import "fmt"

func bubbleSort(nums []int) {
	for i := 0; i < len(nums)-1; i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i] > nums[j] {
				nums[i], nums[j] = nums[j], nums[i]
			}
		}

	}
}

func main() {
	nums := []int{3, 1, 5, 2, 6, 4, 8, 9, 7}
	bubbleSort(nums)
	fmt.Println(nums)
}
