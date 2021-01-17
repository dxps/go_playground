## Order Management Service

### Generating gRPC Code

For generating the Server skeleton and the Client stub based on the protobuf definition that exists in `pbs/order-mgmt.proto` file, you need to make sure the following prerequisites are installed on your system:

- [protoc](https://github.com/protocolbuffers/protobuf/releases/), that is the protocol buffer compiler
    - get the binary release for your system, extract it and have it available in the PATH env var
- gRPC plugin
    - use `go get -u google.golang.org/grpc` (while being outside of a module-based Go project)
- protoc plugin for Go
    - use `go get -u github.com/golang/protobuf/protoc-gen-go` (while being outside of a module-based Go project)

Then you can use the provided `pg-gen.sh` file to generate it or update in case of definition changes. The result (generated code) is stored in `pb/order-mgmt.pb.go` file.

To satisfy the dependencies and be able to just start using the generated code, you'd need to add the protobuf and grpc dependencies to the `go.mod` file, a task that's quickly accomplished by running `go get -u google.golang.org/grpc` while being in this project directory.

### Running the sample

As simple as:
- using `./run_server.sh` to start the server
- using `./run_client.sh` to start the client
