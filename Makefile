.DEFAULT_GOAL := default

revision := $(shell git rev-list -1 HEAD)
version := $(shell git tag | sort -r | head -n1)
netfile := default.nn

ifdef EVALFILE
	netfile := $(EVALFILE)
endif

.PHONY: netgen
netgen:
	rm -f engine/nn.go
	go run -ldflags "-X 'main.netPath=$(netfile)' -X 'main.Version=$(revision)'" netgen/nn.go

build: netgen
	mkdir -p bin
	go build -o bin ./...

ifdef EXE
	mv bin/zahak $(EXE)
endif

run_perft: netgen build
	bin/zahak -perft

run: netgen build
	bin/zahak

test: netgen
	go test ./...

clean:
	go clean ./...
	rm -rf bin

dist:
	echo "Compiling for every OS and Platform"
	mkdir -p bin
	rm -f engine/nn.go
	go run -ldflags "-X 'main.netPath=$(netfile)' -X 'main.Version=$(version)'" netgen/nn.go
	GOOS=linux GOARCH=arm go build -o bin ./... && mv bin/zahak bin/zahak-linux-arm32
	GOOS=linux GOARCH=arm64 go build -o bin ./... && mv bin/zahak bin/zahak-linux-arm64
	GOOS=linux GOARCH=amd64 go build -o bin ./... && mv bin/zahak bin/zahak-linux-amd64
	GOOS=darwin GOARCH=amd64 go build -o bin ./... && mv bin/zahak bin/zahak-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -o bin ./... && mv bin/zahak bin/zahak-darwin-m1-arm64
	GOOS=windows GOARCH=amd64 go build -o bin ./... && mv bin/zahak.exe bin/zahak-windows-amd64.exe
	GOOS=windows GOARCH=386 go build -o bin ./... && mv bin/zahak.exe bin/zahak-windows-386.exe

all: build

default: build
