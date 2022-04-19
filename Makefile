BIN="./bin/previewer"

.PHONY: build
build:
	go build -v -o $(BIN) cmd/previewer/main.go

