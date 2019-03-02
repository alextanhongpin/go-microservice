-include .env
export

install: # Install required modules.
	@go get ./...
	GO111MODULE=on go get ./...

mod: # Initialize go modules and update dependencies.
	GO111MODULE=on go mod init
	GO111MODULE=on go mod tidy 
	GO111MODULE=on go get

start:
	@go run main.go
