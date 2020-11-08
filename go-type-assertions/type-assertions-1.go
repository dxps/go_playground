package main

import "fmt"

func main() {

	var a interface{} = TypeA{"A"}

	// This flavors cause panic:
	// v := a.(TypeB)
	// _ = a.(TypeB)

	// Therefore,
	v, ok := a.(TypeB)
	if ok {
		fmt.Println(" v =", v)
	}

}

// TypeA is the first type to work with.
type TypeA struct {
	Letter string
}

// TypeB is the second type to work with.
type TypeB struct {
	Letter string
}
