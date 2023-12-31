
SHELL := /bin/bash
export

ACCOUNTS_BINARY=accountsApp
BROKER_BINARY=brokerApp
LISTENER_BINARY=listenerApp
USERS_BINARY=usersApp

## up: starts all containers in the background without forcing build
up: build_proto
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_proto build_broker build_accounts build_listener build_users
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

## build_keys: builds the private and public key files for all the neccessary services
build_keys:
	@echo "Building key files..."
	openssl genpkey -algorithm RSA -out broker-service/private_key.pem -pkeyopt rsa_keygen_bits:2048
	openssl rsa -pubout -in broker-service/private_key.pem -out users-service/public_key.pem
	openssl rsa -pubout -in broker-service/private_key.pem -out accounts-service/public_key.pem
	@echo "Done!"

## build_proto: builds the proto grpc files for all the neccessary services
build_proto:
	@echo "Building proto files for user service..."
	cd ./proto && protoc --go_out=../users-service/internals/proto --go_opt=paths=source_relative --go-grpc_out=../users-service/internals/proto --go-grpc_opt=paths=source_relative users.proto
	cd ./proto && protoc --go_out=../broker-service/internals/proto --go_opt=paths=source_relative --go-grpc_out=../broker-service/internals/proto --go-grpc_opt=paths=source_relative users.proto
	@echo "Done!"
	@echo "Building proto files for crypto service..."
	cd ./proto && protoc --go_out=../crypto-service/internals/proto --go_opt=paths=source_relative --go-grpc_out=../crypto-service/internals/proto --go-grpc_opt=paths=source_relative crypto.proto
	cd ./proto && protoc --go_out=../broker-service/internals/proto --go_opt=paths=source_relative --go-grpc_out=../broker-service/internals/proto --go-grpc_opt=paths=source_relative crypto.proto
	@echo "Done!"

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo "Building core binary..."
	cd ./broker-service && go mod download && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./dist/${BROKER_BINARY} ./main.go
	@echo "Done!"

## build_accounts: builds the accounts binary
build_accounts:
	@echo "Building accounts binary..."
	cd ./accounts-service && go mod download && env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./dist/${ACCOUNTS_BINARY} ./main.go
	@echo "Done!"

## build_listener: builds the listener binary
build_listener:
	@echo "Building listener binary..."
	cd ./listener-service && go mod download && env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./dist/${LISTENER_BINARY} ./main.go
	@echo "Done!"

## build_users: builds the users binary
build_users:
	@echo "Building users binary..."
	cd ./users-service && go mod download && env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./dist/${USERS_BINARY} ./main.go
	@echo "Done!"


## ***********************************************************************************

broker:
	cd ./broker-service && go mod download && env SERVICE=.env.local go run main.go

users:
	cd ./users-service && go mod download && env SERVICE=.env.local go run main.go

listener:
	cd ./listener-service && go mod download && env SERVICE=.env.local go run main.go

