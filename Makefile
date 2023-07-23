BINARY=users
OS=linux
TESTS=go test $$(go list ./... | grep -v /vendor/) -cover

#!make
include .env

export ENVIRONMENT=local
export VAULT_ADDR=
export VAULT_TOKEN=

lint:
	golangci-lint run --out-format checkstyle > lint.xml

migrate-gen:
	goose -dir migration create $(name) sql

migrate-up:
	goose -dir migration postgres "host=${DB_HOST} port=${DB_PORT} user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} sslmode=disable" up

migrate-down:
	goose -dir migration postgres "host=${DB_HOST} port=${DB_PORT} user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} sslmode=disable" down
	
local-run:
	go run cmd/main.go

test:
	go test -v ./...