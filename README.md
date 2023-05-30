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
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative auth.proto
```
