package main

import (
	"fmt"

	"github.com/bytecodealliance/wasmtime-go"
)

func main() {

	engine := wasmtime.NewEngine()
	store := wasmtime.NewStore(engine)
	// module, err := wasmtime.NewModuleFromFile(store.Engine, "gcd.wat")
	module, err := wasmtime.NewModuleFromFile(store.Engine, "wasm_rust_hello.wasm")
	check(err)

	modExps := module.Type().Exports()
	if len(modExps) == 0 {
		fmt.Println("No exports exist in the module.")
	} else {
		fmt.Printf("There are %d export(s):\n", len(modExps))
		for _, exp := range modExps {
			fmt.Printf("- %+#v\n", exp.Name())
		}
	}

	modImps := module.Type().Imports()
	if len(modImps) == 0 {
		fmt.Println("No imports exist in the module.")
	} else {
		fmt.Printf("There are %d imports(s):\n", len(modImps))
		for _, imp := range modImps {
			fmt.Printf("- %+#v\n", imp)
		}
	}

	// TODO: To be fixed. This is just for having the compilation.
	// Running it crashes with a seg fault.
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
