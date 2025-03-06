.PHONY: build run dev

SVC_EXEC=go-app_files-mgmt_svc

build:
	@GOARCH=wasm GOOS=js go build -C cmd/pwa -o ../../web/app.wasm
	@go build -C cmd/svc -o ../../$(SVC_EXEC)

start_svc: build
	@./$(SVC_EXEC)

stop_svc:
	@killall -q -HUP $(SVC_EXEC) || echo

dev: stop_svc start_svc
