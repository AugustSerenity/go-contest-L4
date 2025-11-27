run:
	go run ./cmd/main.go

test:
	go test -cover ./...

race-test:
	go test -race ./...

vet:
	go vet ./...

lint:
	golangci-lint run ./...

check: vet lint race-test