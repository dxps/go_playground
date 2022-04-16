## About This Sample

At the time of this writing, TinyGo ver. 0.22 does not yet support Go 1.18+ (it is planned for 0.23 release), therefore:
```shell
❯ tinygo build -o go_run_add.wasm -target wasm ./main.go 
error: requires go version 1.15 through 1.17, got go1.18
❯ 
```
