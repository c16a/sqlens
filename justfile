build-linux-amd64:
    GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o sqlens-linux-amd64 main.go

build-linux-arm64:
    GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o sqlens-linux-arm64 main.go

build-linux-riscv64:
    GOOS=linux GOARCH=riscv64 go build -ldflags="-s -w" -o sqlens-linux-riscv64 main.go

build: build-linux-amd64 build-linux-arm64 build-linux-riscv64