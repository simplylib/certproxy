build: generate
	go build .
generate:
	cd protocol && make
precommit:
	go mod tidy
	go vet -race ./...
	golangci-lint run -v -E gofmt
	gosec -tests ./...
	go test -v -race ./...
	govulncheck -v -test .
