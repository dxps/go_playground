## MirrorFinder

An example of using Go standard libs for creating a Web API.

### Run

Use `go run cmd/main.go` if you want to run it from the shell. <br/>
Otherwise, you can use your IDE of choice.

### Usage

A simple cUrl is enough: `curl -i http://localhost:8000/fastest-mirror`.

You may enhance the output using jq: `curl http://localhost:8000/fastest-mirror | jq`

Examples:
```shell script
$ curl -i http://localhost:8000/fastest-mirror
  HTTP/1.1 200 OK
  Content-Type: application/json
  Date: Sun, 29 Mar 2020 07:04:10 GMT
  Content-Length: 63
  
  {"fastest_url":"http://ftp.hu.debian.org/debian/","latency":76}%
$
```
```shell script
$ curl -s http://localhost:8000/fastest-mirror | jq
  {
    "fastest_url": "http://ftp.ro.debian.org/debian/",
    "latency": 64
  }
$
```
