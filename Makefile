# include variables fomr the .envrc file
include .envrc

# help: print this help message
.PHONY: help
help: 
	@echo 'Usage:'
	@sed -n 's/^#//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# run/api: run the cmd/api application
.PHONY: run/api
run/api:
	go run ./cmd/api -db-dsn=${GREENLIGHT_DB_DSN}?sslmode=disable

# run/api: run the cmd/api application
.PHONY: run/web
run/web:
	go run ./cmd/web

# db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${GREENLIGHT_DB_DSN}

# db/migrations/up apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm 
	@echo 'Running up migrations ....'
	migrate -path ./migrations -database=${GREENLIGHT_DB_DSN}?sslmode=disable up

# db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new: 
	@echo 'Creating migration files for ${name}'
	migrate create -seq -ext=.sql -dir=./migrations ${name}


# =======================#
# QUALITY CONTROL
# =======================#

## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit:
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code ....'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...


## vendor: tidy and vendor depencies
.PHONY: vendor
vendor: 
	@echo 'Tidying and verifying module depencies'
	go mod tidy
	go mod verify
	@echo 'Vendoring depencies...'
	go mod vendor


# =======
# Build
# ======
#

## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	go build -ldflags='-s' -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/api ./cmd/api
