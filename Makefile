DEV_DB_URI = postgres://postgres:1234@localhost:5432/dev-db?sslmode=disable

run: ## Run the application in development mode
	nodemon --watch './**/*.go' --signal SIGTERM --exec APP_ENV=dev 'go' run cmd/app/main.go

dev-db-up: ## Start the development database and Redis
	docker compose up dev-db redis -d

dev-db-rm: ## Remove the development database and Redis containers
	docker compose down -v
	# docker compose rm dev-db -s -f -v

dev-db-migrate-up: ## Run database migrations for the development database
	migrate -path ./internal/repository/migrations -database "$(DEV_DB_URI)" -verbose up

sqlc: ## Generate SQL code with sqlc
	sqlc generate

reset: ## Reset development database, generate SQL code, and run migrations
	$(MAKE) dev-db-rm
	@sleep 10
	$(MAKE) dev-db-up
	@sleep 10
	$(MAKE) sqlc
	$(MAKE) dev-db-migrate-up

help: ## Display this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: run dev-db-up dev-db-rm dev-db-migrate-up sqlc reset help
