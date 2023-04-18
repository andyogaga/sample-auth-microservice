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
	grpcPORT = "50001"
)

func WalletRequestsViaGRPC(service string) (WalletServiceClient, context.Context) {

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", service, grpcPORT), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := NewWalletServiceClient(conn)
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
	RegisterWalletServiceServer(grpcServer, &WalletsServer{})

	log.Printf("gRPC Server started on port :%s", grpcPORT)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve the accounts GRPC server over port: %v", err)
	}
}
