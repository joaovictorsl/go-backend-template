include .env

.EXPORT_ALL_VARIABLES:

.PHONY: migrate-down migrate-up run stop full-stop test

migrate-down:
	goose -dir=migrations postgres ${DATABASE_URL} down

migrate-up:
	goose -dir=migrations postgres ${DATABASE_URL} up

run:
	docker compose up --build -d --remove-orphans
	air

stop:
	docker compose down

full-stop:
	docker compose down -v

args?=./...
test:
	go test ${args}

