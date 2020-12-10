.PHONY: multi install test multi-cross
.PHONY: multi-darwin multi-darwin-386 multi-darwin-amd64
.PHONY: multi-linux multi-linux-386 multi-linux-amd64
.PHONY: multi-windows multi-windows-386 multi-windows-amd64

all: multi

multi:
	go build ./cmd/multi -o ./bin/multi

install:
	go install ./cmd/multi

test:
	go test ./... -cover

multi-cross: multi-darwin multi-linux multi-windows
	@echo "Full cross compilation done."

multi-darwin-386:
	GOOS=darwin GOARCH=386 go build ./cmd/multi -o ./bin/multi-darwin-386
	@echo "Darwin 386 cross compilation done."

multi-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build ./cmd/multi -o ./bin/multi-darwin-amd64
	@echo "Darwin amd64 cross compilation done."

multi-darwin: multi-darwin-386 multi-darwin-amd64
	@echo "Darwin cross compilation done."

multi-linux-386:
	GOOS=linux GOARCH=386 go build ./cmd/multi -o ./bin/multi-linux-386
	@echo "Linux 386 cross compilation done."

multi-linux-amd64:
	GOOS=linux GOARCH=amd64 go build ./cmd/multi -o ./bin/multi-linux-amd64
	@echo "Linux amd64 cross compilation done."

multi-linux: multi-linux-386 multi-linux-amd64
	@echo "Linux cross compilation done."

multi-windows-386:
	GOOS=windows GOARCH=386 go build ./cmd/multi -o ./bin/multi-windows-386
	@echo "Windows 386 cross compilation done."

multi-windows-amd64:
	GOOS=windows GOARCH=amd64 go build ./cmd/multi -o ./bin/multi-windows-amd64
	@echo "Windows amd64 cross compilation done."

multi-windows: multi-windows-386 multi-windows-amd64
	@echo "Windows cross compilation done."
