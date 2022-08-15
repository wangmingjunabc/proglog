.PHONY: compile
compile:
	protoc -I=. --go_out=plugins=grpc,paths=source_relative:. api/v1/log.proto

.PHONY: test
test:
	go test -race ./...

CONFIG_PATH=${HOME}/.proglog/
.PHONY: init
init:
	mkdir -p ${CONFIG_PATH}

.PHONY: gencert
gencert:
	cfssl gencert -initca test/ca-csr.json | cfssljson -bare ca
	cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=test/ca-config.json -profile=server test/server-csr.json | cfssljson -bare server
	mv *.pem *.csr ${CONFIG_PATH}
