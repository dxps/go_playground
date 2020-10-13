package main

import "fmt"

import "unicode/utf8"

func main() {

	var (
		alpha rune = 940
		omega rune = 969
		pi    rune = 960
	)

	// printing the (numeric) values and associated characters
	fmt.Printf(" %6v %6v %6v \n %6[1]c %6[2]c %6[3]c\n", alpha, omega, pi)

	question := "¿Cómo estás?"
	fmt.Println(len(question), "bytes")                    // prints "15 bytes"
	fmt.Println(utf8.RuneCountInString(question), "runes") // prints "12 runes"

	c, size := utf8.DecodeLastRuneInString(question)
	fmt.Printf("First rune: %c %v bytes", c, size)

}
