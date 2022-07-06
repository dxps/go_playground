#!/bin/sh

go build -gcflags="-m -m" -o escape cmd/escape/main.go && rm -f escape

