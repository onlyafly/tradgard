all: build
	heroku local

build:
	go install -race ./cmd/tradgard ./...

deps:
	go get -u github.com/kardianos/govendor
	govendor sync

clean:
	kill -9 $(lsof -i tcp:5000 | grep tradgard | cut -d' ' -f2)

.PHONY: all build clean deps
