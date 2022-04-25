BIN="./bin/previewer"
DOCKER_IMG="previewer:develop"

.PHONY: build
build:
	go build -v -o $(BIN) cmd/previewer/main.go

.PHONY: run
run: build
	$(BIN)

.PHONY: install-linter
install-linter:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.43.0

.PHONY: lint
lint: install-linter
	golangci-lint run ./...

.PHONY: lint-fix
lint-fix:
	golangci-lint run ./... --fix

.PHONY: build-img
build-img:
	docker build \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

.PHONY: up-test
up-test: build-img
	docker-compose \
		--env-file=./deployments/.env \
		-f ./deployments/docker-compose.yaml \
		-f ./deployments/docker-compose.test.yaml \
		up --remove-orphans --force-recreate

.PHONY: test
test:
	go test -race ./internal/...