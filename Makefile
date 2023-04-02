TRANSACTIONS_BINARY=transactionsApp
AUTHENTICATION_BINARY=authenticationApp
BROKER_BINARY=brokerApp
LISTENER_BINARY=listenerApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_broker build_authentication build_listener
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
	cd ./broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ./dist/${BROKER_BINARY} ./main.go
	@echo "Done!"

## build_authentication: builds the authentication binary
build_authentication:
	@echo "Building authentication binary..."
	cd ./authentication-service && env CGO_ENABLED=0 go build -o ./dist/${AUTHENTICATION_BINARY} ./main.go
	@echo "Done!"

## build_listener: builds the listener binary
build_listener:
	@echo "Building listener binary..."
	cd ./listener-service && env CGO_ENABLED=0 go build -o ./dist/${LISTENER_BINARY} ./main.go
	@echo "Done!"

## build_transactions: builds the transactions binary
build_transactions:
	@echo "Building transactions binary..."
	cd ./transactions-service && env CGO_ENABLED=0 go build -o ./dist/${TRANSACTIONS_BINARY} ./main.go
	@echo "Done!"
