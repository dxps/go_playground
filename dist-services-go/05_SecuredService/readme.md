## Writing a Log Package

### Terms

The following terms are being used in the design and code:

| term | description |
| --- | --- |
| `Record` | The data stored in the log |
| `Store` | The file that stores the records. |
| `Index` | The file that stores the index entries. |
| `Segment` | The abstraction that ties a store and an index together. |
| `Log` | The abstraction that ties all the segments together |

### Setup

As prerequisites, the following packages need to be installed, documented in topic sections.

#### Protobufs & gRPC plugin

- Protocol Buffers Compiler (`protoc`)
  ```shell
  $ wget https://github.com/protocolbuffers/protobuf/releases/download/v3.14.0/protoc-3.14.0-linux-x86_64.zip
  $ unzip protoc-3.14.0-linux-x86_64.zip -d ~/apps/protoc
  $ export PATH="$PATH:~/apps/protoc/bin" # add it to your ~/.profile
  ```
- Protocol Buffers Runtime (`gogoprotobuf`)
  `go get github.com/gogo/protobuf/...@v1.3.1`
- gRPC plugin
  `go get google.golang.org/grpc@v1.33.2`

#### CloudFlare's CLI for SSL artifacts

In 05 _SecuredService_ phase, CloudFlare's nice CLI is introduced. `Makefile` is using it for generating the server and client certificates. Install it using:
`go get github.com/cloudflare/cfssl/cmd/cfssl@v1.5.0`

And running `make gencert` should generate the `.csr`, `.pem`, and `-key.pem` files for the CA, server, and two clients (`root` and `nobody`, see `Makefile` for details).

#### Casbin for ACL

Latest `casbin` version 2 was installed (using `go get github.com/casbin/casbin/v2@v2.19.1`) as a dependency.



