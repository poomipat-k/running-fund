APP=./bin/runningApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images"
	docker compose up
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_running
	@echo "Stopping docker images (if running...)"
	docker compose down
	@echo "Building (when required) and starting docker images..."
	docker compose up --build
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker compose down
	@echo "Done!"

## build_running: builds the running app binary as a linux executable
build_running:
	@echo "Building app binary..."
	env GOOS=linux CGO_ENABLED=0 go build -o ${APP} ./cmd/api/
	@echo "Done!"

clear:
	@echo "Deleting unused images..."
	docker image prune -a
	@echo "Done!"

## Migration database

migration_up:
	cd migrations; echo "Inside migrations, Start migration"; \
	goose postgres "host=localhost port=5432 user=poomipat password=running_fund_dev dbname=running_fund_dev sslmode=disable" up
	@echo "Migration Done!"

migration_down:
	cd migrations; echo "Inside migrations, Start migration"; \
	goose postgres "host=localhost port=5432 user=poomipat password=running_fund_dev dbname=running_fund_dev sslmode=disable" down
	@echo "Migration Done!"

migration_status:
	cd migrations; echo "Inside migrations, Start migration"; \
	goose postgres "host=localhost port=5432 user=poomipat password=running_fund_dev dbname=running_fund_dev sslmode=disable" status

test:
	go test -v ./pkg/...