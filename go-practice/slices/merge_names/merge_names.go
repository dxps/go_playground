package main

import "fmt"

func main() {

	// should print Ava, Emma, Olivia, Sophia
	fmt.Println(uniqueNames(
		[]string{"Ava", "Emma", "Olivia"},
		[]string{"Olivia", "Sophia", "Emma"}))
}

func uniqueNames(a, b []string) []string {

	var result []string
	for _, s := range a {
		if !contains(result, s) {
			result = append(result, s)
		}
	}
	for _, s := range b {
		if !contains(result, s) {
			result = append(result, s)
		}
	}
	return result
}

func contains(result []string, item string) bool {

	for _, val := range result {
		if val == item {
			return true
		}
	}
	return false
}
