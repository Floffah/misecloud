set shell := ["bash", "-c"]

default:
	@just --list

build:
	go build -o bin/miseag -ldflags="-s -w" cmd/agent/agent.go
	go build -o bin/miseco -ldflags="-s -w" cmd/controller/controller.go
	go build -o bin/misecl -ldflags="-s -w" cmd/misecl/misecl.go

alias gen := generate
generate:
	go generate ./...

test:
	go test -v ./...

lint:
	golangci-lint run

alias fmt := format
format:
	go fmt ./...
	golangci-lint fmt