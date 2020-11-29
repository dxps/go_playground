
# This variable specifies where to put the generated certs.
CONFIG_PATH=${HOME}/.dist-services-go



# An init target to creat that config directory.
.PHONY: init
init:
	mkdir -p ${CONFIG_PATH}



.PHONY: gencert
gencert:
	cfssl gencert \
		-initca test/ca-csr.json | cfssljson -bare ca

	# the CA
	cfssl gencert \
		-ca=ca.pem \
		-ca-key=ca-key.pem \
		-config=test/ca-config.json \
		-profile=server \
		test/server-csr.json | cfssljson -bare server

	# the client
	cfssl gencert \
		-ca=ca.pem \
		-ca-key=ca-key.pem \
		-config=test/ca-config.json \
		-profile=client \
		test/client-csr.json | cfssljson -bare client

	# for the multi client test: a `root` (super user) and a `nobody` user

	cfssl gencert \
		-ca=ca.pem \
		-ca-key=ca-key.pem \
		-config=test/ca-config.json \
		-profile=client \
		-cn="root" \
		test/client-csr.json | cfssljson -bare root-client
	
	cfssl gencert \
		-ca=ca.pem \
		-ca-key=ca-key.pem \
		-config=test/ca-config.json \
		-profile=client \
		-cn="nobody" \
		test/client-csr.json | cfssljson -bare nobody-client
	
	mv *.pem *.csr ${CONFIG_PATH}



$(CONFIG_PATH)/model.conf:
	cp test/model.conf $(CONFIG_PATH)/model.conf

$(CONFIG_PATH)/policy.csv:
	cp test/policy.csv $(CONFIG_PATH)/policy.csv



.PHONY: test
test:
test: $(CONFIG_PATH)/policy.csv $(CONFIG_PATH)/model.conf
	go test -race ./...



# Generating the Go files from protobuf definitions.
.PHONY: compile
compile:
	protoc api/v1/*.proto \
		--gogo_out=\
Mgogoproto/gogo.proto=github.com/gogo/protobuf/proto,plugins=grpc:. \
		--proto_path=\
$$(go list -f '{{ .Dir }}' -m github.com/gogo/protobuf) \
		--proto_path=.

