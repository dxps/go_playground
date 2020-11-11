package main

import "fmt"

// CounterFactory returns a closure that will
func CounterFactory(j int) func() int {
	i := j
	// This is a closure that is returned: it's the function bundled together
	// with references ot its surrounding state (the lexical environment).
	// It has access to the outer function's scope from this inner function.
	return func() int {
		i++
		return i
	}
}

func main() {
	incr := CounterFactory(1)
	fmt.Printf(" %d\n", incr())
	fmt.Printf(" %d\n", incr())
	fmt.Printf(" %d\n", incr())
}
