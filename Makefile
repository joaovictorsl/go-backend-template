include .env

.PHONY: migration-down migration-up run

migrate-down:
	migrate -path=migrations -database=${DATABASE_URL} down

migrate-up:
	migrate -path=migrations -database=${DATABASE_URL} up

run:
	docker compose up --build -d --remove-orphans
	air
