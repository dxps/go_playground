## A hello test using `httprouter`

The initial purpose of this sample is to evaluate performance (although might not be so relevant for such a simple use case), compared with a similar [Rust based hello test using Rocket](https://github.com/dxps/rust_playground/tree/master/rust_rocket_hello) sample.

### Run

Use the classic `go run ./cmd/server/main.go` to run it.

<br/>

### Usage

`curl -v localhost:8001` to access the `/` path.

For a small stress test, you can use `stress.sh` script.

As mentioned before, the implementation is too simple to be relevant enough as comparison between Go and Rust. It was made simply out of curiosity. Here are some figures snapshot when sending 10k requests over 50 connections:
| values | Go  | Rust |
| --- | --- | --- |
| total  | 0.09294186s    | 0.11211325s   |
| stats  | avg request time: 0.0139ms <br/>median: 0ms <br/> 95th percentile: 0ms <br/> 99th percentile: 1ms| avg request time: 0.0128ms <br/> median: 0ms <br/> 95th percentile: 0ms <br/> 99th percentile: 1ms |

<br/>
