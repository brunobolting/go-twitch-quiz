all: build

SHELL := env ENV=$(ENV) $(SHELL)
ENV ?= dev

BIN_DIR = $(PWD)/bin

clean:
	rm -rf bin/*

dependencies:
	go mod download

build: dependencies build-app

build-app:
	go build -o ./bin main.go

run:
	go run .

test:
	go test -tags testing ./...
