package main

var global *T

type T struct {
	i int
}

func NewT() *T {
	return &T{10}
}

func main() {
	t := NewT() // new instance of T
	global = t  // have it referred from `global`
}
