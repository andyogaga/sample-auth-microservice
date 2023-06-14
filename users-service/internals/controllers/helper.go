package controller

import (
	"crypto/rsa"
	"fmt"
	"log"
	"net"
	"os"

	events "users-service/internals/event"
	proto "users-service/internals/proto"

	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
)

func SetupGRPCRequestsListener(messageQueueConfig *events.Config) {
	listenAddr := fmt.Sprintf("localhost:%s", grpcPORT)
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(LoggerInterceptor(messageQueueConfig), AuthorizationInterceptor))
	proto.RegisterUserServiceServer(grpcServer, &UsersServer{})

	log.Printf("gRPC Server started on: %s", listenAddr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve the users GRPC server over port: %v", err)
	}
}

func loadPublicKey() (*rsa.PublicKey, error) {
	publicKeyStr := os.Getenv("RSA_PUBLIC_KEY")
	publicKeyBytes := []byte(publicKeyStr)
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)

	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	return publicKey, nil
}
