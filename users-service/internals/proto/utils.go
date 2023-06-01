package requests

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpcPORT = "50002"
)

func UserRequestsViaGRPC(service string) (UserServiceClient, context.Context) {

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", service, grpcPORT), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return c, ctx
}

func SetupGRPCRequestsListener() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	RegisterUserServiceServer(grpcServer, &UsersServer{})

	log.Printf("gRPC Server started on port :%s", grpcPORT)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve the users GRPC server over port: %v", err)
	}
}
