package main

// This function is imported from 'external world'.
func add(x, y int) int

func main() {

	// That imported function is being used here.
	println("adding two numbers: ", add(2, 3))

}
