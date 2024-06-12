run-develop: 
	@go run cmd/main.go

build:
	@go build -o bin/gotext cmd/main.go

db-up:
	@cd ./config && docker compose up db -d

db-stop: 
	@cd ./config && docker compose down db

docker-up:
	@cd ./config && docker compose up

docker-down:
	@cd ./config && docker compose down
