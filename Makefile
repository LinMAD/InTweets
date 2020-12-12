export GO111MODULE=on

all: vet lint

vet:
	@go vet -all ./

lint:
	@golint -set_exit_status ./...

build: # Simple one
	@go build .
