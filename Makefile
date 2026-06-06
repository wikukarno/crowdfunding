.PHONY: run build test cover vet wire docs tidy migrate-up migrate-down

run:
	go run .

build:
	go build -o bin/server .

test:
	go test -race ./...

cover:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

vet:
	go vet ./...

# Regenerate the Wire injector after changing providers.
wire:
	go run github.com/google/wire/cmd/wire ./...

# Regenerate Swagger docs after changing handler annotations.
docs:
	go run github.com/swaggo/swag/cmd/swag init --parseDependency --parseInternal -g main.go

tidy:
	go mod tidy

# Requires golang-migrate and a DATABASE_URL, e.g.
#   export DATABASE_URL="postgres://postgres:postgres@127.0.0.1:5432/crowdfunding?sslmode=disable"
migrate-up:
	migrate -path ./migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path ./migrations -database "$(DATABASE_URL)" down
