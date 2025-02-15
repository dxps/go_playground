build:
	GOARCH=wasm GOOS=js go build -C src -o ../web/app.wasm
	go build -C src -o ../nav

run: build
	./nav

stop:
	killall -q -HUP nav || echo

dev: stop run
	

.PHONY: build run dev

