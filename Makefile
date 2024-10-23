
run-api:
	go run cmd/main.go api

run-reader:
	go run cmd/main.go reader

build:

up:
	docker-compose up -d

down:
	docker-compose down

install:
	echo "=== Installing dependencies ==="
	go mod tidy
	go mod vendor
	echo "Done"