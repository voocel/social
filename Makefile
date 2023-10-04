VERSION:=$(shell grep 'VERSION' pkg/version/version.go | awk '{ print $$4 }' | tr -d '"')

default:gen

protofile ?= *.proto

gen:
	@protoc -I. --go_out=./protos \
	 	    --go-grpc_out=./protos \
	 	   ./protos/$(protofile)

pb:
	@protoc -I="./protos" --go_out=./protos \
			--go-grpc_out=./protos \
			./protos/gate.proto ./protos/node.proto ./protos/message.proto;

## build: Compile a program into an executable file
build:
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags '-w -s' -o social

## compress: Compress executable files
compress: build
	upx --brute social

## ents: Ent generate
ents:
	@ent generate ./ent/schema --target ./ent

test:
	go test

## install-protobuf: Install protobuf plugins
install-protobuf:
	go install \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
		google.golang.org/protobuf/cmd/protoc-gen-go \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc

## help: Show this help screen
help: Makefile
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo '' | awk '{printf "\033[36m%-30s\n", $$1}'
	@echo 'Available targets are:'
	@echo ''
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo ''
	@echo 'make gen protofile=[your proto filename]'
	@echo ''

## version: Display social version
version:
	@echo ${VERSION}

## bench: Benchmarks are run sequentially
bench:
	go test -benchmem -cpuprofile cpu.out -memprofile mem.out -run=^$ github.com/voocel/social/benchmark -bench ^Benchmark$
	go tool pprof -svg bench.test cpu.out > cpu.svg
	go tool pprof -svg bench.test mem.out > mem.svg