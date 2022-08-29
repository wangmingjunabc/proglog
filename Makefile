.PHONY: compile

$(CONFIG_PATH)/model.conf:
	cp test/model.conf $(CONFIG_PATH)/model.conf
$(CONFIG_PATH)/policy.csv:
	cp test/policy.csv $(CONFIG_PATH)/policy.csv

compile:
	protoc -I=. --go_out=plugins=grpc,paths=source_relative:. api/v1/log.proto

.PHONY: test
test: $(CONFIG_PATH)/model.conf $(CONFIG_PATH)/policy.csv
	go test -race ./...

CONFIG_PATH=${HOME}/.proglog/
.PHONY: init
init:
	mkdir -p ${CONFIG_PATH}

.PHONY: gencert
gencert:
	cfssl gencert -initca test/ca-csr.json | cfssljson -bare ca
	cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=test/ca-config.json -profile=server test/server-csr.json | cfssljson -bare server
	cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=test/ca-config.json -profile=client test/client-csr.json | cfssljson -bare client
	cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=test/ca-config.json -profile=client test/client-csr.json | cfssljson -bare client
	cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=test/ca-config.json -profile=client -cn="root" test/client-csr.json | cfssljson -bare root-client
	cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=test/ca-config.json -profile=client -cn="nobody" test/client-csr.json | cfssljson -bare nobody-client
	mv *.pem *.csr ${CONFIG_PATH}

TAG ?= 0.0.1
build-docker:
	 docker buildx build  --platform linux/amd64  -t github.com/wangmingjunabc/proglog:$(TAG) .