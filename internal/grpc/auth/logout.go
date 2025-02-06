package auth

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	sl "sso/pkg/logger"
	"sso/protos/gen/go/sso"
)

func (s *serverAPI) Logout(ctx context.Context, in *sso.LogoutRequest) (*sso.LogoutResponse, error) {
	const op = "auth.Logout"

	refreshToken, err := s.getRefreshTokenInMetadata(ctx)
	if err != nil {
		return &sso.LogoutResponse{Success: false}, status.Error(codes.Internal, "failed to logout")
	}

	if err = s.auth.RevokeTokens(ctx, refreshToken); err != nil {
		sl.Log.Error(op, "failed to revoke tokens", sl.Err(err))
		return &sso.LogoutResponse{Success: false}, status.Error(codes.Internal, "failed to logout")
	}

	sl.Log.Info(op, "User logged out successfully")
	return &sso.LogoutResponse{Success: true}, nil
}
