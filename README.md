# kendi-server

server for kendi

## Architecture

- This is an attempt to implement a clean architecture, in case you don’t know it yet, here’s a reference https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html

## Requirements/dependencies

- Docker
- Docker-compose

## Getting Started

- Starting API in development mode

```sh
docker-compose up
```

- Stopping API in development mode

```sh
docker-compose down
```

## Convert .Proto files

- Install Protobuf

```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

- Download protobuf file from `https://grpc.io/docs/protoc-installation/` and add the `protoc` binary to your $PATH (Preferrably `/go/bin`)

- `cd` to the paths of the `.proto` file and run

```sh
cd proto
protoc --go_out=../users-service/internals/proto --go_opt=paths=source_relative --go-grpc_out=../users-service/internals/proto --go-grpc_opt=paths=source_relative users.proto
protoc --go_out=../broker-service/internals/proto --go_opt=paths=source_relative --go-grpc_out=../broker-service/internals/proto --go-grpc_opt=paths=source_relative users.proto
```

### Local dev env setup

## Setting up Rabbitmq

1. Install rabbitmq [https://www.rabbitmq.com/] and use homebrew
2. Start Rabbitmq service locally `brew services run rabbitmq`
3. Local management plugin at `http://localhost:15672`
4. Stop Rabbitmq service locally `brew services stop rabbitmq`

## Setting up Postgres

1. Install rabbitmq [https://www.postgresql.org/download/] and use homebrew
2. Start Rabbitmq service locally `brew services run postgresql@14`
3. Local management can be done using pgAdmin
4. Stop Rabbitmq service locally `brew services stop postgresql@14`
5. Having issues? Restart with `brew services restart postgresql`

Use `brew services list` to see your running homebrew services

## Running servers locally

1. In different Terminal tabs, run these commands.

```sh
make broker
make users
make listener
```

## P2P provider

`https://docs-api.cryptoprocessing.io/#wallet-object`
