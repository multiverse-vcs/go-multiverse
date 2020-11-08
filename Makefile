.PHONY: all

CGO_ENABLED=1

# release:
# 	GOOS=darwin GOARCH=386 go build -o bin/multi-darwin-386 main.go
# 	GOOS=darwin GOARCH=amd64 go build -o bin/multi-darwin-amd64 main.go
# 	GOOS=freebsd GOARCH=386 go build -o bin/multi-freebsd-386 main.go
# 	GOOS=freebsd GOARCH=amd64 go build -o bin/multi-freebsd-amd64 main.go
# 	GOOS=freebsd GOARCH=arm go build -o bin/multi-freebsd-arm main.go
# 	GOOS=linux GOARCH=386 go build -o bin/multi-linux-386 main.go
# 	GOOS=linux GOARCH=amd64 go build -o bin/multi-linux-amd64 main.go
# 	GOOS=linux GOARCH=arm go build -o bin/multi-linux-arm main.go
# 	GOOS=linux GOARCH=arm64 go build -o bin/multi-linux-arm64 main.go
# 	GOOS=openbsd GOARCH=386 go build -o bin/multi-openbsd-386 main.go
# 	GOOS=openbsd GOARCH=amd64 go build -o bin/multi-openbsd-amd64 main.go
# 	GOOS=openbsd GOARCH=arm go build -o bin/multi-openbsd-arm main.go
# 	GOOS=windows GOARCH=386 go build -o bin/multi-windows-386 main.go
# 	GOOS=windows GOARCH=amd64 go build -o bin/multi-windows-amd64 main.go

all:
	go build -o bin/multi main.go