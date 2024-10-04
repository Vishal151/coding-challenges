.PHONY: build run-lb run-be clean

build:
	go build -o bin/lb cmd/lb/main.go
	go build -o bin/be cmd/be/main.go

run-lb: build
	./bin/lb

run-be: build
	./bin/be

clean:
	rm -rf bin