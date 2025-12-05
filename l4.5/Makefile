PORT ?= 8080

run:
	@echo "Запуск сервера на порту $(PORT)..."
	PORT=$(PORT) go run ./cmd/main.go

test:
	go test -cover ./...

race-test:
	go test -race ./...

vet:
	go vet ./...

lint:
	golangci-lint run ./...

check: vet lint race-test