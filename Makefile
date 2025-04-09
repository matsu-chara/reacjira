NAME=reacjira
BUILD_VERSION=latest

build:
	go build

lint:
	go tool golangci-lint run --enable gofmt

test:
	go test ./...

fmt:
	gofmt -s -w .

clean:
	go clean
	rm -f ${NAME}

run:
	REACJIRA_CONFIG_NAME=config.secret.toml go run main.go

docker:
	docker build --platform linux/amd64 -t ${NAME}:$(BUILD_VERSION) .

release: docker
	docker push ${NAME}:$(BUILD_VERSION)
