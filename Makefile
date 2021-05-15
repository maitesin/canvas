tools-generate:
	go get -u github.com/matryer/moq

tools-lint: tools-generate
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

generate:
	go generate ./...

test: generate
	go test -cover -v ./...

lint: generate
	golangci-lint run

run:
	go run cmd/canvas/main.go
