## Go-based gRPC > Greet basic example

### Prereqs

The followings were installed from a location outside of this project, since we don't want to add them explictly as dependencies into this project.

- Protocol Buffers Compiler (protoc)<br/>
  ```shell
  $ curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v3.14.0/protoc-3.14.0-linux-x86_64.zip
  $ unzip protoc-3.14.0-linux-x86_64.zip -d ~/apps/protoc
  $ rm protoc-3.14.0-linux-x86_64.zip
  $ export PATH=$PATH:${HOME}/apps/protoc/bin   # add it to your ~/.profile
  ```
- gRPC-Go (the gRPC plugin for Protocol Buffers) using `go get -u google.golang.org/grpc`
- protobuf (Go support for Protocol Buffers) using `go get -u github.com/golang/protobuf/protoc-gen-go`

### 

