.PHONY: build
build:
	mkdir -p ./bin
	go build -o ./bin/go2rs ./cmd/go2rs
