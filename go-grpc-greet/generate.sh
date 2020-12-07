#!/bin/sh

protoc greet/greetpb/greet.proto --go_out=plugins=grpc:.
