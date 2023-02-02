include .envrc

## audit: tidy dependencies, format and vet all code
.PHONY: audit
audit:
	@echo 'Tidying and verifying the dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...

## run: run the cmd/api application
.PHONY: run
run:
	go run ./cmd/api -db-dsn=${PARKOUR_DB_DSN}

## build: build the binary for the application
.PHONY: build
build:
	@echo 'Building the binaries'
	go build -ldflags='-s' -o=./bin/dev/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/prod/api ./cmd/api

## psql: open the database
.PHONY: psql
psql:
	psql ${PARKOUR_DB_DSN}

## migration: create new migration files
.PHONY: migration
migration:
	@echo 'Creating migration files for ${table}'
	migrate create -seq -ext=.sql -dir=./migrations ${table}

## db_up: run up migrations
.PHONY: db_up
db_up:
	@echo "Runnig up migrations..."
	migrate -path ./migrations -database ${PARKOUR_DB_DSN} up

## db_down: run down migrations
.PHONY: db_down
db_down:
	@echo "Runnig down migrations..."
	migrate -path ./migrations -database ${PARKOUR_DB_DSN} down