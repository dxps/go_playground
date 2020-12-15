
# (Re)Generate the gRPC Code
gen-grpc:
	protoc --go_out=pb --go-grpc_out=pb     \
		--go_opt=paths=source_relative      \
		--go-grpc_opt=paths=source_relative \
		proto/*.proto
