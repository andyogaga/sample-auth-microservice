
ACCOUNTS_BINARY=accountsApp
BROKER_BINARY=brokerApp
LISTENER_BINARY=listenerApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_broker build_accounts build_listener
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo "Building core binary..."
	cd ./broker-service && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./dist/${BROKER_BINARY} ./main.go
	@echo "Done!"

## build_accounts: builds the accounts binary
build_accounts:
	@echo "Building accounts binary..."
	cd ./accounts-service && env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./dist/${ACCOUNTS_BINARY} ./main.go
	@echo "Done!"

## build_listener: builds the listener binary
build_listener:
	@echo "Building listener binary..."
	cd ./listener-service && env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./dist/${LISTENER_BINARY} ./main.go
	@echo "Done!"

