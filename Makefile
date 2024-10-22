.PHONY: build run

build:
	mkdir -p bin
	go build -o ./bin/doubtfire ./cmd/doubtfire

run: build
	./bin/doubtfire
