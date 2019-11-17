package main

import (
	"fmt"
	"os"
)

const filename = "functions.txt"

// printer1 is an example of returning multiple parameters.
func printer1(msg string) (string, error) {
	msg += "\n"
	_, err := fmt.Printf("printer1 > %s", msg)
	return msg, err
}

// printer2 is an example of using defer keyword,
// declaring the code that is called at the end of the function's execution.
func printer2(msg string) error {
	f, err := os.Create(filename)
	fmt.Printf("printer2 > message %q was saved to %s file\n", msg, filename)
	defer f.Close()
	defer fmt.Printf("printer2 > closing %s file.\n", f.Name())
	// doing some logic ...
	return err
}

// printer3 shows how to use named return paramters.
func printer3() (e error) {
	fmt.Printf("printer3 > Removing %s file.\n", filename)
	e = os.Remove(filename)
	return
}

// printer4 is a variadic function
// (allowing one or multiple parameters of the specified type).
func printer4(format string, msgs ...string) {
	for _, msg := range msgs {
		fmt.Printf("printer4 > "+format+"\n", msg)
	}
}

func main() {
	msg := "Hello, World!"
	appendedMessage, err := printer1(msg)
	if err == nil {
		fmt.Printf("main > printer1 returned %q\n", appendedMessage)
	}
	printer2(msg)
	printer3()
	printer4("%s", msg, msg)
}
