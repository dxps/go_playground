package main

import "fmt"

// This sample code shows the two forms of the `for range`:
// - Value Semantic Form
//      - this is `for i, val := range myArray {...}`
//      - note that `val` is a copy of the value that is in the array at that index
// - Pointer Semantic Form
//      - this is `for i := range myArray {...}`

func main() {

	// Using the value semantic form of the `for range`.
	friends := [3]string{"Annie", "Bobby", "Charley"}
	fmt.Printf("[value semantic form]\tbefore: %s", friends[1])
	for i, v := range friends {
		if i == 1 {
			friends[i] = "Joe"
			fmt.Printf("\tafter: %s\n", v)
		}
	}

	// Using the pointer semantic form of the `for range`.
	friends = [3]string{"Annie", "Bobby", "Charley"}
	fmt.Printf("[pointer semantic form]\tbefore: %s", friends[1])
	for i := range friends {
		if i == 1 {
			friends[i] = "Joe"
			fmt.Printf("\tafter: %s\n", friends[1])
		}
	}

}
