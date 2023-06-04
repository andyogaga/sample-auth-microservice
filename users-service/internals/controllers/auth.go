package controller

import (
	context "context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"users-service/internals/dto"
	proto "users-service/internals/proto"
	services "users-service/internals/services"
)

const (
	grpcPORT = "50002"
)

var userService services.UserService

func SetupService(_userService *services.UserService) {
	userService = *_userService
}

type UsersServer struct {
	proto.UnimplementedUserServiceServer
}

func UserRequestsViaGRPC(service string) (proto.UserServiceClient, context.Context) {

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", service, grpcPORT), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := proto.NewUserServiceClient(conn)
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
	proto.RegisterUserServiceServer(grpcServer, &UsersServer{})

	log.Printf("gRPC Server started on port :%s", grpcPORT)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve the users GRPC server over port: %v", err)
	}
}

func (c *UsersServer) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	fmt.Println(req)
	user, err := userService.Register(dto.CreateUser{Phone: req.Phone, Country: req.Country})
	if err != nil {
		return nil, err
	}
	return &proto.CreateUserResponse{Message: user.UserId}, nil
}
