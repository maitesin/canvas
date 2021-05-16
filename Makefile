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

start-infra:
	docker-compose up -d

stop-infra:
	docker-compose stop

remove-infra:
	docker-compose down

run: start-infra
	go run cmd/canvas/main.go
