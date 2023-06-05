package controller

import (
	context "context"
	"log"

	events "users-service/internals/event"

	"google.golang.org/grpc"
)

func RequestLoggerInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
	messageQueueConfig events.Config,
) (resp interface{}, err error) {
	log.Printf("Received gRPC request: %s", info.FullMethod)
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
