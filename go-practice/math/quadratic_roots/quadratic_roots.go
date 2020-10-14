package main

import (
	"fmt"
	"math"
)

func findRoots(a, b, c float64) (float64, float64) {

	x1 := (-b + math.Sqrt(b*b-4*a*c)) / (2 * a)
	x2 := (-b - math.Sqrt(b*b-4*a*c)) / (2 * a)
	return x1, x2
}

func main() {

	x1, x2 := findRoots(2, 10, 8)
	fmt.Printf("Roots: %f, %f", x1, x2)
}
