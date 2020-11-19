package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {

	x := make([]string, 2, 3)
	x[0] = "g"
	x[1] = "o"

	y := append(x, "t")

	fmt.Printf(" x = %s   ", x)
	printSliceInternals(x)

	fmt.Printf(" y = %s ", y)
	printSliceInternals(y)
}

func printSliceInternals(s []string) {
	h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	fmt.Printf(" internals:  Len=%v  Cap=%v  Data=%v\n", h.Len, h.Cap, h.Data)
}
