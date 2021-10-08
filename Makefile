.DEFAULT_GOAL := default

ifneq ($(OS), Windows_NT)
	revision := $(shell git rev-list -1 HEAD || echo dev)
	version := $(shell git tag | sort -r | head -n1)
endif

netfile := default.nn

ifdef EVALFILE
	netfile := $(EVALFILE)
endif

RM=rm -f engine/nn.go
MKDIR=mkdir -p bin
MV=mv bin/zahak $(EXE)
ifeq ($(OS), Windows_NT)
	RM=del engine\nn.go
	MKDIR=IF not exist bin (mkdir bin)
	MV=move bin\zahak.exe $(EXE).exe
endif

.PHONY: netgen
netgen: clean
	go run -ldflags "-X 'main.netPath=$(netfile)' -X 'main.Version=$(revision)'" netgen/nn.go

build: netgen
	$(MKDIR)
	go build -o bin ./...

ifdef EXE
	$(MV)
endif

run_perft: netgen build
	bin/zahak -perft

run: netgen build
	bin/zahak

test: netgen
	go test ./...

clean:
	go clean ./...
	$(RM)

dist: clean
	$(MKDIR)
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
