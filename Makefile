ARGS ?=

run:
	go run ./cmd/hexlet-path-size $(ARGS)
build:
	go build -o bin/hexlet-path-size ./cmd/hexlet-path-size