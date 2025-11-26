.PHONY: build run

build:
	go build -o mycut cmd/main.go
	chmod +x test.sh

run: build
	./test.sh
