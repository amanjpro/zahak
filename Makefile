.DEFAULT_GOAL := default

ifneq ($(OS), Windows_NT)
	revision := $(shell git rev-list -1 HEAD || echo dev)
	version := $(shell git tag | sort -r | head -n1)
endif

netfile := default.nn

ifdef EVALFILE
	netfile := $(EVALFILE)
endif

RM=rm -f engine/nn.go engine/nn_*.go

MKDIR=mkdir -p bin
MV=mv bin/zahak $(EXE)
FLAGS=CC=cc CGO_ENABLED="1"
ifeq ($(OS), Windows_NT)
	RM=del engine\nn.go engine\nn_*.go
	MKDIR=IF not exist bin (mkdir bin)
	MV=move bin\zahak.exe $(EXE).exe
	FLAGS=
endif

.PHONY: netgen
netgen: clean
	$(FLAGS) go run -gcflags "-B" -ldflags "-X 'main.netPath=$(netfile)' -X 'main.Version=$(revision)'" netgen/nn.go

tcec:
	$(FLAGS) go run -gcflags "-B" -ldflags "-X 'main.netPath=$(netfile)' -X 'main.Version=$(TCEC_VERSION)'" netgen/nn.go
	$(MKDIR)
	$(FLAGS) go build -gcflags "-B" --ldflags '-linkmode external -extldflags "-static"' -o bin ./...
	mv bin/zahak bin/zahak-linux-amd64-$(TCEC_VERSION)-avx

build: netgen
	$(MKDIR)
	$(FLAGS) go build -gcflags "-B" -o bin ./...

ifdef EXE
	$(MV)
endif

debug: netgen
	$(MKDIR)
	$(FLAGS) go build -o bin ./...
	mv bin/zahak bin/zahak_debug

run_perft: netgen build
	bin/zahak -perft

run: netgen build
	bin/zahak

debug_run: netgen debug
	bin/zahak_debug

test: netgen
	go test ./...

clean:
	go clean ./...
	$(RM)

cross-build: clean
	$(MKDIR)
	echo "!!!! WARNING !!!! Cross build will not support Syzygy Probing"
	$(FLAGS) go run -ldflags "-X 'main.netPath=$(netfile)' -X 'main.Version=$(version)'" netgen/nn.go
	$(FLAGS) GOOS=linux GOARCH=arm go build -gcflags "-B" -o bin ./... && mv bin/zahak bin/zahak-linux-arm32
	$(FLAGS) GOOS=linux GOARCH=arm64 go build -gcflags "-B" -o bin ./... && mv bin/zahak bin/zahak-linux-arm64
	$(FLAGS) GOOS=linux GOARCH=amd64 go build -gcflags "-B" -o bin ./... && mv bin/zahak bin/zahak-linux-amd64
	$(FLAGS) GOOS=darwin GOARCH=amd64 go build -gcflags "-B" -o bin ./... && mv bin/zahak bin/zahak-darwin-amd64
	$(FLAGS) GOOS=darwin GOARCH=arm64 go build -gcflags "-B" -o bin ./... && mv bin/zahak bin/zahak-darwin-m1-arm64
	$(FLAGS) GOOS=windows GOARCH=amd64 go build -gcflags "-B" -o bin ./... && mv bin/zahak.exe bin/zahak-windows-amd64.exe
	$(FLAGS) GOOS=windows GOARCH=386 go build -gcflags "-B" -o bin ./... && mv bin/zahak.exe bin/zahak-windows-386.exe

all: build

default: build
