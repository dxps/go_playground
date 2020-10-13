
## Run

Start the server using `go run cmd/server/main.go` or the provided `run_server.sh` script.

## Usage

Add log records using:

```shell
$ curl -X POST localhost:8080 -d '{ "record": { "value": "12345678"}}'
$ curl -X POST localhost:8080 -d '{ "record": { "value": "1234567890123456" } }'
$ curl -X POST localhost:8080 -d '{ "record": { "value": "TGV0J3MgR28gIzEK" } }'
$ curl -X POST localhost:8080 -d '{ "record": { "value": "TGV0J3MgR28gIzIK" } }'
```

Note that the value must be a base64 encoded string because that is what [encoding/json](https://golang.org/pkg/encoding/json/#Marshal) expects and does under the hoods for `[]byte` values.

Each such request gets in response the offset where the record has been saved. Internally, this works just like an index.

Get the log records using:

```shell
$ curl -X GET localhost:8080 -d '{"offset":0}'
$ curl -X GET localhost:8080 -d '{"offset":1}'
```
