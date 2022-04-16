package main

import (
	"fmt"

	"github.com/bytecodealliance/wasmtime-go"
)

func main() {

	engine := wasmtime.NewEngine()
	store := wasmtime.NewStore(engine)
	module, err := wasmtime.NewModuleFromFile(store.Engine, "wasm_rust_hello.wasm")
	// module, err := wasmtime.NewModuleFromFile(store.Engine, "gcd.wat")
	check(err)

	modExps := module.Type().Exports()
	if len(modExps) == 0 {
		fmt.Println("No exports found in the module.")
	} else {
		fmt.Printf("There are %d export(s):\n", len(modExps))
		for _, exp := range modExps {
			fmt.Printf("- %+#v\n", exp.Name())
		}
	}

	imports := make([]wasmtime.AsExtern, 4)
	instance, err := wasmtime.NewInstance(store, module, imports)
	check(err)

	// After we've instantiated we can get the exported `add` function and call it.
	add := instance.GetExport(store, "print_hello").Func()
	if add == nil {
		panic("'print_hello' is not a function")
	}

	val, err := add.Call(store)
	check(err)
	fmt.Printf("print_hello result: %s\n", val.(string))
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
