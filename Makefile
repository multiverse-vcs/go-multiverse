release:
	GOOS=freebsd GOARCH=386 go build -o bin/multi-freebsd-386
	GOOS=freebsd GOARCH=amd64 go build -o bin/multi-freebsd-amd64
	GOOS=darwin GOARCH=386 go build -o bin/multi-darwin-386
	GOOS=darwin GOARCH=amd64 go build -o bin/multi-darwin-amd64
	GOOS=linux GOARCH=386 go build -o bin/multi-linux-386
	GOOS=linux GOARCH=amd64 go build -o bin/multi-linux-amd64
	GOOS=windows GOARCH=386 go build -o bin/multi-windows-386
	GOOS=windows GOARCH=amd64 go build -o bin/multi-windows-amd64

all: build

build:
	go build -o bin/multi

install: build
	cp bin/multi /usr/local/bin/

.PHONY: build release install
