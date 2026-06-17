.PHONY: run migrate-up migrate-down migrate-down-all tidy docker-dev docker-prod

run:
	go run ./cmd/server

migrate-up:
	go run ./cmd/migrate -direction up

migrate-down:
	go run ./cmd/migrate -direction down

migrate-down-all:
	go run ./cmd/migrate -direction down-all

tidy:
	go mod tidy

docker-db:
	docker compose up db -d

docker-dev:
	docker compose up --build

docker-prod:
	docker compose -f compose.prod.yml up --build -d
