package main

import (
	"fmt"
	"sort"
)

func containsDuplicateMap(nums []int) bool {
	m := make(map[int] int)
	for _,num := range nums {
		m[num]++
		if(m[num] > 1) {
			return true
		}
	}
	return false
}

func containsDuplicateSort(nums []int) bool {
	clone := make([]int, len(nums))
	copy(clone, nums)
	sort.Ints(clone)
	current := nums[0]
	for i := 1; i< len(nums); i++ {
		if(nums[i] == current) {
			return true
		} else {
			current = nums[i]
		}
	}
	return false
}

func main1() {
	nums := []int{1,2,3,4,5,6,7,8,9,10,11,12,1}
	fmt.Println(containsDuplicateMap(nums))
	fmt.Println(containsDuplicateSort(nums))
}
