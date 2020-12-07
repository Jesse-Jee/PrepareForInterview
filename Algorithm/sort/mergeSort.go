package main

import "fmt"

func mergeSort(nums []int, left, right int) {
	if left >= right {
		return
	}

	mid := (left + right)>>1

	mergeSort(nums, left, mid)
	mergeSort(nums, mid+1, right)

	merge(nums, left, mid, right)

}

func merge(nums []int, left, mid, right int) {
	tmp := make([]int, right-left+1)
	i, j, k := left, mid+1, 0

	for i <= mid && j <= right {
		if nums[i] <= nums[j] {
			tmp[k] = nums[i]
			i++
		} else {
			tmp[k] = nums[j]
			j++
		}
		k++
	}

	for i <= mid {
		tmp[k] = nums[i]
		i++
		k++
	}

	for j <= right {
		tmp[k] = nums[j]
		j++
		k++
	}

	for p := 0; p < len(tmp); p++ {
		nums[left+p] = tmp[p]
	}

}

func main() {
	nums := []int{3, 1, 5, 2, 6, 4, 8, 9, 7}
	mergeSort(nums, 0, len(nums)-1)
	fmt.Println(nums)
}
