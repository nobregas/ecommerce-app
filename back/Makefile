build:
	@go build -o bin/ecommerce-app-backend cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/ecommerce-app-backend

migrate:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

migrate-force:
	@go run cmd/migrate/main.go force $(version)

database-up:
	@docker-compose up -d

database-down:
	@docker-compose down