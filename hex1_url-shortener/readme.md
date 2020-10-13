# Hexagonal Architecture :: example 1 :: URL Shortener

Notes:
- The main goal was to evaluating a project structure that aims to follow Hexagonal (Clean) Architecture
- This is Tensor's example presented in _Building Hexagonal Microservices with Go_ series on YouTube. See [Part 1](https://www.youtube.com/watch?v=rQnTtQZGpg8).
- Some naming adjustments have been made: for example, I like `ShortUrl` model name instead of `Redirect` as it sounds 
  more related to the domain concept (and it shouldn't be related to the action of redirecting as outcome) 

## Run

### Prereqs

If you want to use a local Redis instance, you can simply start it in a Docker container 
using `docker run --name redis -p 6379:6379 -d redis`

You may run `docker logs redis` to see if Redis container successfully started.

### Start

Use `DB_TYPE=redis DB_URL=redis://localhost:6379 go run cmd/main.go` to start the server.

### Stop

Currently, it listens for interrupt signals that you can send by 
doing a `CTRL+C` in the terminal where the server was started.

## Usage

See below the two main HTTP calls that can be used.

After each call is processed, you can see in the terminal processing related log entries that include timestamp, source IP, request, and processing time.

### Adding a new entry

```shell script
$ curl -X POST localhost:8080 -d '{ "url":  "https://devisions.org" }'
{"code":"qcIl1tNMR","url":"https://devisions.org","created_at":1598090181}
$
```

### Getting an entry based on the code

```shell script
$ curl -v localhost:8080/qcIl1tNMR
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 8080 (#0)
> GET /qcIl1tNMR HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.64.1
> Accept: */*
>
< HTTP/1.1 301 Moved Permanently
< Content-Type: text/html; charset=utf-8
< Location: https://devisions.org
< Date: Sat, 22 Aug 2020 09:58:30 GMT
< Content-Length: 56
<
<a href="https://devisions.org">Moved Permanently</a>.

* Connection #0 to host localhost left intact
* Closing connection 0
$ 
```

This request, done through a Browser, will redirect (see the status code of 301) to the value of the _Location_ header.
