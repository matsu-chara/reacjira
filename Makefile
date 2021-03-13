clean:
	go clean
	rm -f reacjira

lint:
	docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:v1.38.0 golangci-lint run ./...

build:
	go build

test:
	go test ./...

fmt:
	gofmt -s -w .

run:
	REACJIRA_CONFIG_NAME=config.secret.toml go run main.go
