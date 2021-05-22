revision := $(shell git rev-list -1 HEAD)
version := $(shell git tag | sort -r | head -n1)

build:
	mkdir -p bin
	go build -ldflags "-X 'main.version=$(revision)'" -o bin ./...
ifdef EXE
	mv bin/zahak bin/$(EXE)
endif

run_perft:
	mkdir -p bin
	go build -o bin ./...
	bin/zahak -perft

run:
	mkdir -p bin
	go build -ldflags "-X 'main.version=$(revision)'" -o bin ./...
	bin/zahak

test:
	go test ./...

clean:
	go clean ./...
	rm -rf bin

dist:
	echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=arm go build -ldflags "-X 'main.version=$(version)'" -o bin ./... && mv bin/zahak bin/zahak-linux-arm32
	GOOS=linux GOARCH=arm64 go build -ldflags "-X 'main.version=$(version)'" -o bin ./... && mv bin/zahak bin/zahak-linux-arm64
	GOOS=linux GOARCH=amd64 go build -ldflags "-X 'main.version=${version}'" -o bin ./... && mv bin/zahak bin/zahak-linux-amd64
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X 'main.version=${version}'" -o bin ./... && mv bin/zahak bin/zahak-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -ldflags "-X 'main.version=${version}'" -o bin ./... && mv bin/zahak bin/zahak-darwin-m1-arm64
	GOOS=windows GOARCH=amd64 go build -ldflags "-X 'main.version=${version}'" -o bin ./... && mv bin/zahak.exe bin/zahak-windows-amd64.exe
	GOOS=windows GOARCH=386 go build -ldflags "-X 'main.version=${version}'" -o bin ./... && mv bin/zahak.exe bin/zahak-windows-386.exe

all: build
