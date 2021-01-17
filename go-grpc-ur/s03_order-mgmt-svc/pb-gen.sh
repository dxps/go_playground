#!/bin/sh

## Make the destination directory, if not exists.
mkdir -p pb

protoc -I pbs pbs/order-mgmt.proto --go_out=plugins=grpc:pb

