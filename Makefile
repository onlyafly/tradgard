PID = /tmp/tradgard.pid

all: serve

deps:
	go get -u github.com/kardianos/govendor
	govendor sync
	brew install fswatch

serve: restart
	@fswatch -o . | xargs -n1 -I{}  make restart

before_restart:
	@echo "Restarting..."

kill:
	@pkill -P `cat $(PID)` || true
	@killall tradgard || true

build:
	@go install -race ./cmd/tradgard ./...

restart: before_restart kill build
	@heroku local & echo $$! > $(PID)

.PHONY: all build deps kill serve restart
