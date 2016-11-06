all: build
	heroku local

build:
	go install -race ./cmd/tradgard ./...

deps:
	go get -u github.com/kardianos/govendor
	govendor sync

.PHONY: all build deps
