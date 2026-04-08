run:
	go run ./cmd/main.go

test:
	go test -v ./...

test-file:
	@if "$(FILE)"=="" ( \
		echo Please provide FILE=path\to\file_test.go & exit /b 1 \
	)
	cmd /C "set DEV=true && go test -v $(FILE)"

all: 
	run test

tidy:
	go mod tidy