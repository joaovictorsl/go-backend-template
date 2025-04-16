include .env

.PHONY: migration-down migration-up run

GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}

migration-down:
	GOOSE_DRIVER=${GOOSE_DRIVER} \
	GOOSE_DBSTRING=${GOOSE_DBSTRING} \
	goose -dir=migrations down

migration-up:
	GOOSE_DRIVER=${GOOSE_DRIVER} \
	GOOSE_DBSTRING=${GOOSE_DBSTRING} \
	goose -dir=migrations up

run:
	docker compose up --build
