all: build

SHELL := env ENV=$(ENV) $(SHELL)
ENV ?= dev

BIN_DIR = $(PWD)/bin

clean:
	rm -rf bin/*

dependencies:
	go mod download

build: dependencies build-fixture build-app

build-app:
	go build -o ./bin/app main.go

build-fixture:
	go build -o ./bin/fixture fixture/fixture.go

run:
	go run .

test:
	go test -tags testing ./...
