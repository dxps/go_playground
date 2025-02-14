build:
	GOARCH=wasm GOOS=js go build -o web/app.wasm
	go build -o hello-world

run: build
	./hello-world

.PHONY: build run