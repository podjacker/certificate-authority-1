default: insecure services test

secure:
	go generate ./vendor/github.com/go-ocf/kit/security
.PHONY: secure

insecure:
	OCF_INSECURE=TRUE go generate ./vendor/github.com/go-ocf/kit/security
.PHONY: insecure

test:
	go test ./service ./...
.PHONY: test

test_docker:
	docker-compose run --rm golang -e DOCKER=1 make test
.PHONY: test_docker

services:
	docker-compose pull zookeeper kafka mongo
	docker-compose up -d zookeeper kafka mongo
.PHONY: services

stop:
	docker-compose down
.PHONY: stop

proto/get:
	go get github.com/gogo/protobuf/proto
	go get github.com/gogo/protobuf/protoc-gen-gogofaster
	go get github.com/gogo/protobuf/gogoproto
	export PATH=${PATH}:${GOPATH}/bin
.PHONY: proto/get

proto/generate:
	protoc -I=. -I=${GOPATH}/src -I=${GOPATH}/src/github.com/gogo/protobuf/protobuf --gogofaster_out=${GOPATH}/src pb/authorization.proto
	protoc -I=. -I=${GOPATH}/src -I=${GOPATH}/src/github.com/gogo/protobuf/protobuf --gogofaster_out=${GOPATH}/src pb/cert.proto
	protoc -I=. -I=${GOPATH}/src -I=${GOPATH}/src/github.com/gogo/protobuf/protobuf --go_out=plugins=grpc:${GOPATH}/src pb/service.proto
.PHONY: proto/generate


deps:
	go get -v -t ./...
.PHONE: deps

build: deps
	go build
.PHONY: build

publish-coverage:
	curl -s https://codecov.io/bash > .codecov && \
	chmod +x .codecov && \
	./.codecov -f coverage.txt
.PHONY: publish-coverage

cover:
	go generate ./...
	go test ./service -race -covermode=atomic -coverprofile=coverage.txt
.PHONY: cover

clean:
	@find . -name \.covdecov -type f -delete
	@find . -name \.coverage.txt -type f -delete
	@rm -f gover.coverage
.PHONY: clean



