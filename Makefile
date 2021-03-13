ci: clean devdeps lint build test

clean:
	go clean
	rm -f reacjira

devdeps:
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

lint:
	golangci-lint run

build:
	go build

test:
	go test ./...

fmt:
	gofmt -s -w .

run:
	REACJIRA_CONFIG_NAME=config.secret.toml go run main.go
