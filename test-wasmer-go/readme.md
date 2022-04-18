## Sample 1 - Hello

This is WIP.

Currently, due to WASM packaging (see [details](./s1_hello/wasm_rust_hello/readme.md)), there are some imports that are needed, thus `go run main.go` crashes with:
```
panic: Cannot get module instance: Missing import: `wasi_snapshot_preview1`.`fd_write`
```

Note that `wasmtime --invoke print_hello wasm_rust_hello.wasm` runs fine.
