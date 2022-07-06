package main

type T struct {
	i int
}

func NewT() *T {
	return &T{10}
}

func main() {
	t := NewT() // new instance of T
	_ = t       // just to make compiler happy
}
