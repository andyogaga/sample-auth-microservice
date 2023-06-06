package controller

import (
	context "context"
	"fmt"
	"log"
	"net"

	events "users-service/internals/event"
	proto "users-service/internals/proto"

	"google.golang.org/grpc"
)

func requestLoggerInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
	messageQueueConfig *events.Config,
) (resp interface{}, err error) {
	payload := events.Payload{
		Name: events.EVENT,
		Data: struct {
			Method string
			Body   interface{}
		}{
			Method: info.FullMethod,
			Body:   req,
		},
	}
	messageQueueConfig.LogEventViaRabbit(&payload)
	return handler(ctx, req)
}

func SetupGRPCRequestsListener(messageQueueConfig *events.Config) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		return requestLoggerInterceptor(ctx, req, info, handler, messageQueueConfig)
	}))
	proto.RegisterUserServiceServer(grpcServer, &UsersServer{})

	log.Printf("gRPC Server started on port :%s", grpcPORT)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve the users GRPC server over port: %v", err)
	}
}
