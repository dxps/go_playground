#!/bin/sh

go build -gcflags="-m -m" -o no_escape cmd/no_escape/main.go && rm -f no_escape

