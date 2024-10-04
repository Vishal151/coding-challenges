.PHONY: build run-lb run-be test docker-build docker-run clean

build:
	go build -o bin/lb cmd/lb/main.go
	go build -o bin/be cmd/be/main.go

run-lb: build
	./bin/lb

run-be: build
	./bin/be

test:
	go test ./...

docker-build:
	docker-compose build

docker-run:
	docker-compose up

clean:
	rm -rf bin