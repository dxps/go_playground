package main

import (
	"fmt"
	"io/ioutil"

	wasmer "github.com/wasmerio/wasmer-go/wasmer"
)

func main() {

	wasmBytes, _ := ioutil.ReadFile("wasm_rust_hello.wasm")

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
	wasmFn, err := instance.Exports.GetFunction("print_hello")
	if err != nil {
		panic(fmt.Sprintf("Cannot get exported function: %v", err))
	}

	// Calls that exported function with Go standard values. The WebAssembly
	// types are inferred and values are casted automatically.
	result, _ := wasmFn()

	fmt.Println(result)

}
