package auth

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	sl "sso/pkg/logger"
	"sso/protos/gen/go/sso"
)

func (s *serverAPI) RefreshToken(ctx context.Context, in *sso.RefreshTokenRequest) (*sso.RefreshTokenResponse, error) {
	const op = "auth.RefreshToken"

	refreshToken, err := s.getRefreshTokenInMetadata(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	accessToken, newRefreshToken, err := s.auth.RefreshTokens(ctx, refreshToken)
	if err != nil {
		sl.Log.Error(op, "failed creating access token", sl.Err(err))
		return nil, status.Error(codes.Internal, "internal server error")
	}

	outgoingMD := metadata.Pairs("refresh_token", newRefreshToken)
	grpc.SendHeader(ctx, outgoingMD)

	return &sso.RefreshTokenResponse{AccessToken: accessToken}, nil
}
