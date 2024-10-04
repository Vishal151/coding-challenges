.PHONY: build run-lb run-be test docker-build docker-run clean generate-cert run-all run-lb-with-ui

build:
	go build -o bin/lb cmd/lb/main.go
	go build -o bin/be cmd/be/main.go

run-lb: build
	./bin/lb

run-be1: build
	./bin/be 8081

run-be2: build
	./bin/be 8082

run-all: build
	@echo "Stopping any running instances..."
	-pkill -f "./bin/be" || true
	-pkill -f "./bin/lb" || true
	@echo "Starting backend servers..."
	./bin/be 8081 & ./bin/be 8082 & 
	@echo "Waiting for backends to start..."
	sleep 2
	@echo "Starting load balancer..."
	BACKEND_URLS=http://localhost:8081,http://localhost:8082 ./bin/lb

run-lb-with-ui: build
	@echo "Starting load balancer with web UI..."
	@BACKEND_URLS=http://localhost:8081,http://localhost:8082 ./bin/lb

test:
	go test ./...

docker-build:
	docker-compose build

docker-run:
	docker-compose up

clean:
	rm -rf bin
	docker-compose down

generate-cert:
	openssl req -x509 -newkey rsa:4096 -keyout server.key -out server.crt -days 365 -nodes -subj "/CN=localhost"

build-loadtest:
	go build -o bin/loadtest cmd/loadtest/main.go

run-loadtest: build-loadtest
	./bin/loadtest -url http://localhost:8080 -c 10 -n 1000

run-lb-round-robin: build
	@echo "Starting load balancer with round-robin algorithm..."
	@BACKEND_URLS=http://localhost:8081,http://localhost:8082 ./bin/lb -algorithm round-robin

run-lb-least-connections: build
	@echo "Starting load balancer with least connections algorithm..."
	@BACKEND_URLS=http://localhost:8081,http://localhost:8082 ./bin/lb -algorithm least-connections

run-lb-ip-hash: build
	@echo "Starting load balancer with IP hash algorithm..."
	@BACKEND_URLS=http://localhost:8081,http://localhost:8082 ./bin/lb -algorithm ip-hash