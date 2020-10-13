package main

import (
	"fmt"
)

func main() {
	fruits := make([]string, 3)
	fruits[0] = "Apple"
	fruits[1] = "Banana"
	fruits[2] = "Grape"

	inspectSlice(fruits)
	appendSlice(&fruits, "Orange")
	// Note in this second 'inspection' that the new value force the creation
	// of a new array underneath with a double capacity (6 instead of 3).
	inspectSlice(fruits)
}

func inspectSlice(slice []string) {
	fmt.Printf("Length:%d Capacity:%d\n", len(slice), cap(slice))
	for i, s := range slice {
		fmt.Printf("at idx:%d ref(pointer) addr:%p value:%s\n", i, &slice[i], s)
	}
}

// `appendSlice` appends the provided value to the provided slice.
// As in Go everything is passed by value, it's using pointers to reflect the change
// on the underlying array structure.
func appendSlice(slice *[]string, value string) {
	*slice = append(*slice, value)
}
