package main

import (
	"fmt"
	"math/rand"
)

func main() {

	num := rand.Intn(10) + 1
	fmt.Printf("num = %v\n", num)
	num = rand.Intn(10) + 1
	fmt.Println("num =", num)

}
