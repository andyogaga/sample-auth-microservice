package controller

import (
	context "context"
	"crypto/rsa"
	"fmt"
	"log"
	events "users-service/internals/event"

	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ContextToken string

const (
	BROKER_TOKEN ContextToken = "token"
)

func VerifyToken(tokenString string, publicKey *rsa.PublicKey) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing error")
		}
		return publicKey, nil
	})
	if err != nil {
		return err
	}
	_, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func AuthorizationInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("error getting token from grpc")
	}
	token := md["token"][0]
	publicKey, err := loadPublicKey()
	if err != nil {
		log.Print("error is", err.Error())
		return nil, fmt.Errorf("error loading public key")
	}
	// Verify the authorization token using the public key
	verifyError := VerifyToken(token, publicKey)
	if verifyError != nil {
		return nil, verifyError
	}

	return handler(ctx, req)
}

func LoggerInterceptor(messageQueueConfig *events.Config) grpc.UnaryServerInterceptor {

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
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
}
