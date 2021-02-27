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

dist:
	echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=arm go build -o bin ./... && mv bin/{zahak,zahak-linux-arm}
	GOOS=linux GOARCH=amd64 go build -o bin ./... && mv bin/{zahak,zahak-linux-amd64}
	GOOS=darwin GOARCH=amd64 go build -o bin ./... && mv bin/{zahak,zahak-darwin-amd64}
	GOOS=windows GOARCH=amd64 go build -o bin ./... && mv bin/{zahak.exe,zahak-windows-amd64.exe}
	GOOS=windows GOARCH=386 go build -o bin ./... && mv bin/{zahak.exe,zahak-windows-386.exe}

all: build
