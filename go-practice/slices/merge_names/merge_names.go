package main

import "fmt"

func main() {
	// should print Ava, Emma, Olivia, Sophia
	fmt.Println(uniqueNames(
		[]string{"Ava", "Emma", "Olivia"},
		[]string{"Olivia", "Sophia", "Emma"}))
}

func uniqueNames(a, b []string) []string {
	var r []string
	r = collectUniqueValues(r, a...)
	r = collectUniqueValues(r, b...)
	return r
}

func collectUniqueValues(r []string, vals ...string) []string {
	for _, val := range vals {
		if !contains(r, val) {
			r = append(r, val)
		}
	}
	return r
}

func contains(result []string, item string) bool {
	for _, val := range result {
		if val == item {
			return true
		}
	}
	return false
}
