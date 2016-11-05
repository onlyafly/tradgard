all:
	go run cmd/tradgard/main.go

deps:
	go get -u github.com/kardianos/govendor
	govendor sync

.PHONY: all deps
