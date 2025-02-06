package auth

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	sl "sso/pkg/logger"
)

func (s *serverAPI) getRefreshTokenInMetadata(ctx context.Context) (string, error) {
	op := "sso.getRefreshTokenInMetadata"
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md["refresh_token"]) == 0 {
		sl.Log.Warn(op, "refresh token is missing in metadata")
		return "", status.Error(codes.InvalidArgument, "refresh token is required")
	}

	refreshToken := md["refresh_token"][0]
	if refreshToken == "" {
		sl.Log.Warn(op, "received empty refresh token")
		return "", status.Error(codes.InvalidArgument, "refresh token is required")
	}
	return refreshToken, nil
}
