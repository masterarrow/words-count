build:
	@go build -ldflags="-s -w" -o wc ./cmd

test:
	@go test -v ./...

vet:
	@go vet ./...
	@golangci-lint run

race:
	@go run -race ./cmd
