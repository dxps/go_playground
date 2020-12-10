#!/bin/sh

# protoc -I=greet/pb --go_out=plugins=grpc:greet/pb greet.proto

protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    greet/pb/greet.proto
