PID = /tmp/tradgard.pid

include .env

all: serve

deps:
	go get -u github.com/mattes/migrate
	dep ensure
	brew install fswatch

serve: restart
	@fswatch -o . | xargs -n1 -I{}  make restart

before_restart:
	@echo "Restarting..."

kill:
	@pkill -P `cat $(PID)` || true
	@killall tradgard || true

migrate-version:
	migrate -url $(DATABASE_URL) -path ./etc/db version

migrate-down-1:
	migrate -url $(DATABASE_URL) -path ./etc/db migrate -1

migrate-up:
	migrate -url $(DATABASE_URL) -path ./etc/db up

build:
	@go install -race ./cmd/tradgard ./...

restart: before_restart kill build
	@heroku local & echo $$! > $(PID)

.PHONY: all build deps kill migrate-down-1 migrate-up migrate-version serve restart
