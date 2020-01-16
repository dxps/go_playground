package main

import (
	"fmt"
	"log"
)

const (
	pluginName    = "TestPlugin"
	pluginVersion = 0x00010000
)

// Load is loading the plugin.
func Load(register func(name string, version uint64) error) error {
	err := register(pluginName, pluginVersion)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	log.Printf("%s > loaded version: 0x%x\n", pluginName, pluginVersion)
	return nil
}

// Unload is unloading the plugin.
func Unload() error {
	fmt.Printf("%s > unloaded version: 0x%x\n", pluginName, pluginVersion)
	return nil
}

// Test is a simple usage example/feature.
func Test(data string) string {
	return "hello " + data
}
