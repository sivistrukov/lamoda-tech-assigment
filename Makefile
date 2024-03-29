all: test up

.PHONY: swag
swag:
	swag init -g cmd/main.go

.PHONY: migration
migration:
	migrate create -ext sql -dir ./migrations -seq $(name)

.PHONY: test
test:
	go test -v -timeout 30s ./...

.PHONY: down
down:
	docker compose down

.PHONY: up
up:
	docker compose up --build

