VERSION:=$(shell grep 'VERSION' pkg/version/version.go | awk '{ print $$4 }' | tr -d '"')

default:gen

gen:
	@protoc -I. --go_out=./proto \
	 	    --go-grpc_out=./proto \
	 	   ./proto/$(protofile)

build:
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags '-w -s' -o social

compress: build
	upx --brute social

test:
	go test

install-protobuf:
	go install \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
		google.golang.org/protobuf/cmd/protoc-gen-go \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc

help:
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@echo '    help               Show this help screen'
	@echo '    build              Compile a program into an executable file'
	@echo '    compress           Compress executable files'
	@echo '    install-protobuf   Install protobuf plugins'
	@echo '    version            Display social version'
	@echo ''
	@echo 'make gen proto=[your proto filename]'
	@echo ''

version:
	@echo ${VERSION}