package auth

import (
	"context"
	"google.golang.org/grpc"
	"sso/protos/gen/go/sso"
)

type Auth interface {
	SignIn(ctx context.Context, username string, password string, appID int) (string, string, error)
	RefreshTokens(ctx context.Context, refreshToken string) (string, string, error)
	RevokeTokens(ctx context.Context, refreshToken string) error
}

type serverAPI struct {
	sso.UnimplementedAuthServer
	auth Auth
}

func Register(gRPCServer *grpc.Server, auth Auth) {
	sso.RegisterAuthServer(gRPCServer, &serverAPI{auth: auth})
}
