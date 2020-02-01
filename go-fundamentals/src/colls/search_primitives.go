package main

import (
	"fmt"
	"sort"
)

var numbers = []int{3, 2, 5, 7, 8}

func main() {
	fmt.Println(" numbers =", numbers)
	fmt.Println(" Is numbers sorted in asc order?", sort.IntsAreSorted(numbers))

	indexOf5 := sort.SearchInts(numbers, 5)
	fmt.Print(" Searching for 5 in numbers. Result (index):", indexOf5)
	if indexOf5 < len(numbers) {
		fmt.Printf(" Found! numbers[%d] = %d\n", indexOf5, numbers[indexOf5])
	}

	indexOf8 := sort.SearchInts(numbers, 8)
	fmt.Print(" Searching for 8 in numbers. Result (index):", indexOf8)
	if indexOf8 < len(numbers) {
		fmt.Printf(" Found! numbers[%d] = %d\n", indexOf8, numbers[indexOf8])
	}

	indexOf9 := sort.SearchInts(numbers, 9)
	fmt.Print(" Searching for 9 in numbers. Result (index):", indexOf9)
	if indexOf9 < len(numbers) {
		fmt.Printf(" Found! numbers[%d] = %d\n", indexOf9, numbers[indexOf9])
	} else {
		fmt.Println(" Not found!")
	}
}
