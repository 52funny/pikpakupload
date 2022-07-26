.PHONY: all build linux check run

BIN=pikpakupload

all: check build

build:
	go build -o ${BIN} .

linux:
    CGO_ENABLED=0  GOOS=linux  GOARCH=amd64  go build -o ${BIN} .

check:
	go mod tidy
	go fmt $(go list ./... | grep -v /vendor/)

run:
	go run .

