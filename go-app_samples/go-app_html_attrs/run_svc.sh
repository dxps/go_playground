#!/bin/sh

## npx nodemon --exec "go run" ./main.go --signal SIGTERM
##
## Using fresh tool
## (install it using `go install github.com/gravityblast/fresh@latest`)
##

## Not used since I cannot tell him to watch `.` and run `go build cmd/gmd`
## fresh -c ./ops/run_dev.conf

## Using `entr` (from https://eradman.com/entrproject/)
find . -name "*.go" | entr -r make dev
