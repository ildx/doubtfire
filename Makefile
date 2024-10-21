.PHONY: all build run clean

all: build

build:
	go build -o bin/doubtfire ./cmd/doubtfire

run:
	go run ./cmd/doubtfire $(ARGS)

clean:
	rm -rf bin