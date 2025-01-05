DEV_DB_URI = postgres://postgres:1234@localhost:5432/dev-db?sslmode=disable

run:
	nodemon --watch './**/*.go' --signal SIGTERM --exec APP_ENV=dev 'go' run main.go
dev-db-up:
	docker compose up dev-db redis -d

dev-db-rm:
	docker compose down -v
	# docker compose rm dev-db -s -f -v


dev-db-migrate-up:
	migrate -path ./internal/db/migration -database "$(DEV_DB_URI)" -verbose up

.PHONY: run dev-db-up dev-db-rm db-migrate-up