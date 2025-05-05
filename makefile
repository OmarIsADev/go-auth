BINARY_NAME=app

all: deps build run

deps:
	go mod download tidy

build:
	go build -o $(BINARY_NAME) ./cmd/main.go

run:
	./$(BINARY_NAME)

rm:
	rm -rf $(BINARY_NAME)