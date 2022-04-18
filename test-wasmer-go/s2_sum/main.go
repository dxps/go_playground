package main

import (
	"fmt"
	"io/ioutil"

	wasmer "github.com/wasmerio/wasmer-go/wasmer"
)

func main() {

	wasmBytes, _ := ioutil.ReadFile("wasm_rust_sum_build.wasm")

	engine := wasmer.NewEngine()
	store := wasmer.NewStore(engine)

	// Compiles the module
	module, _ := wasmer.NewModule(store, wasmBytes)

	// Instantiates the module
	importObj := wasmer.NewImportObject()
	instance, err := wasmer.NewInstance(module, importObj)
	if err != nil {
		panic(fmt.Sprintf("Cannot get module instance: %v", err))
	}

	// Gets the `sum` exported function from the WebAssembly instance.
	sum, err := instance.Exports.GetFunction("sum")
	if err != nil {
		panic(fmt.Sprintf("Cannot get exported function: %v", err))
	}

	// Calls that exported function with Go standard values. The WebAssembly
	// types are inferred and values are casted automatically.
	result, _ := sum(5, 37)

	fmt.Println(result) // 42!

}
