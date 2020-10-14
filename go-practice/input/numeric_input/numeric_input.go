package main

import (
	"fmt"
	"strconv"
)

type UserInput interface {
	Add(rune)
	GetValue() string
}

type NumericInput struct {
	input string
}

func (i *NumericInput) Add(r rune) {

	v := string(r)
	_, err := strconv.Atoi(v)
	if err != nil {
		return
	}
	i.input = i.input + v
}

func (i *NumericInput) GetValue() string {
	return i.input
}

func main() {

	var input UserInput = &NumericInput{}
	input.Add('1')
	input.Add('a')
	input.Add('0')
	fmt.Println(input.GetValue())
}
