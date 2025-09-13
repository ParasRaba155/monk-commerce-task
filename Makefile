fmt:
	go fmt ./...

build:
	go build -o bin/app cmd/*.go

start:build
	./bin/app
