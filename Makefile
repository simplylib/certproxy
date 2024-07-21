.PHONY: build
build:
	go build .

.PHONY: run-server
run-server:
	go run . server

.PHONY: precommit
precommit:
	-go mod tidy
	-go vet ./...
	-golangci-lint run -v -E gofmt
	-gosec -tests ./...
	-go test -v -race ./...
	-govulncheck -v -test .
