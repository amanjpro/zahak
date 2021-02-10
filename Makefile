build:
	mkdir -p bin
	go build -o bin ./...

run_perft:
	mkdir -p bin
	go build -o bin ./...
	bin/zahak -perft

run:
	mkdir -p bin
	go build -o bin ./...
	bin/zahak

test:
	go test ./...

clean:
	go clean ./...
	rm -rf bin

# compile:
# 	echo "Compiling for every OS and Platform"
# 	GOOS=linux GOARCH=arm go build -o bin/main-linux-arm main.go
# 	GOOS=linux GOARCH=arm64 go build -o bin/main-linux-arm64 main.go
# 	GOOS=freebsd GOARCH=386 go build -o bin/main-freebsd-386 main.go

all: build
