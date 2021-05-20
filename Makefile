tools:
	go install github.com/matryer/moq@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

generate:
	go generate ./...

test: generate
	go test -cover -v ./...

test-integration: generate
	go test -tags integration -v ./...

lint: generate
	golangci-lint run

start-infra:
	docker-compose up -d

stop-infra:
	docker-compose stop

remove-infra:
	docker-compose down

run:
	go run cmd/canvas/main.go
