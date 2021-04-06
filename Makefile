revision := $(shell git rev-list -1 HEAD)
version := $(shell git tag | sort -r | head -n1)

build:
	mkdir -p bin
	go build -ldflags "-X 'main.version=$(revision)'" -o bin ./...

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
	GOOS=windows GOARCH=amd64 go build -ldflags "-X 'main.version=${version}'" -o bin ./... && mv bin/zahak.exe bin/zahak-windows-amd64.exe
	GOOS=windows GOARCH=386 go build -ldflags "-X 'main.version=${version}'" -o bin ./... && mv bin/zahak.exe bin/zahak-windows-386.exe

all: build
