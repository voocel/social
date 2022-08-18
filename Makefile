default:gen

gen:
	@protoc -I. --go_out=./proto \
	 	    --go-grpc_out=./proto \
	 	   ./proto/$(protofile)

test:
	go test

help:
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ' make gen proto=[your proto filename]'
	@echo ''