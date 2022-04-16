## About This Sample

A minimal hello world example, built as a WASI target (WASM).

<br/>

### Building

- prerequisite: `cargo install cargo-wasi`
- build: `cargo wasi build`
   - The generated WASM file is `target/wasm32-wasi/debug/wasm_rust_hello.wasm`.
- execute: `wasmtime --invoke print_hello wasm_rust_hello.wasm`
