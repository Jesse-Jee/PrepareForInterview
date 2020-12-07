package main

import "fmt"

func selectSort(nums []int) {

	for i := 0; i < len(nums)-1; i++ {
		min := i
		for j := i + 1; j < len(nums); j++ {
			if nums[j] < nums[min] {
				min = j
			}
		}
		nums[i], nums[min] = nums[min], nums[i]
	}

}

func main() {
	nums := []int{3, 9, 8, 5, 6, 4, 2, 7, 1}
	selectSort(nums)
	fmt.Println(nums)
}
