ARGS ?=

run:
	go run ./cmd/hexlet-path-size $(ARGS)
build:
	go build -o bin/hexlet-path-size ./cmd/hexlet-path-size
lint:
	golangci-lint run
lint-fix:
	golangci-lint run --fix
test:
	go test -v ./... $(ARGS)