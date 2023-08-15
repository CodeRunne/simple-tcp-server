init:
	go mod init github.com/codrunne/tcp_server

tidy:
	go mod tidy

run:
	go run main.go

build:
	go build -o bin/main main.go