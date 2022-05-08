BIN="./bin/previewer"
DOCKER_IMG="previewer:develop"

build:
	go build -v -o $(BIN) cmd/previewer/main.go

run: build
	$(BIN)

install-linter:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.45.2

lint: install-linter
	golangci-lint run ./...

lint-fix:
	golangci-lint run ./... --fix

build-img:
	docker build \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

up-test: build-img
	docker-compose \
		--env-file=./deployments/.env \
		-f ./deployments/docker-compose.yaml \
		-f ./deployments/docker-compose.test.yaml \
		up --remove-orphans --force-recreate

test:
	go test -race ./internal/...

.PHONY: build run install-linter lint lint-fix build-img up-test test