package main

import (
	"fmt"
)

// printer takes an array named w as input (a copy of it)
func printer1(w [4]string) {
	fmt.Print("printer1 > ")
	for _, word := range w {
		fmt.Printf("%s ", word)
	}
	fmt.Println()
	w[2] = "blue"
}

// printer2 taks a slice as input.
// a slice is a window into an array.
func printer2(w []string) {
	fmt.Print("printer2 > ")
	for _, word := range w {
		fmt.Printf("%s ", word)
	}
	fmt.Println()
	w[2] = "blue"
}

func main() {

	// an array is passed as copy to functions
	wordsArray := [4]string{"the", "quick", "brown", "fox"}
	// a slice is passed as reference to functions
	wordsSlice := []string{"the", "quick", "brown", "fox"}
	printer1(wordsArray)
	printer1(wordsArray)
	printer2(wordsSlice)
	printer2(wordsSlice)

	// making a slice of strings (up to 4 items) and populate it later
	words := make([]string, 0, 4)
	fmt.Printf("words slice has length:%d and capacity:%d\n", len(words), cap(words))
	words = append(words, "the")
	words = append(words, "quick")
	words = append(words, "brown")
	words = append(words, "fox")
	fmt.Printf("words slice has length:%d and capacity:%d\n", len(words), cap(words))
	// this append will double the capacity
	words = append(words, "jumps")
	fmt.Printf("words slice has length:%d and capacity:%d\n", len(words), cap(words))

}
