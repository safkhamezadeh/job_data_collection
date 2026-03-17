run:
	go run ./cmd/main.go

test:
	go test -v ./...

all: 
	run test

tidy:
	go mod tidy