package main

import (
	"fmt"
	"github.com/letiantech/hotplugin"
)

func main() {
	options := hotplugin.ManagerOptions{
		Dir:    "./plugins",
		Suffix: ".so",
	}
	hotplugin.StartManager(options)
	result := hotplugin.Call("TestPlugin", "Test", "my world")
	fmt.Println(result...)
}
