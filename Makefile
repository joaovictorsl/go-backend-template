include .env

.PHONY: migration-down migration-up run

DATABASE_URL=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable

migrate-down:
	migrate -path=migrations -database=${DATABASE_URL} down

migrate-up:
	migrate -path=migrations -database=${DATABASE_URL} up

run:
	docker compose up --build
