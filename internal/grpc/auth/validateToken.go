package auth

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sso/internal/service"
	"sso/protos/gen/go/sso"
)

func (s *serverAPI) ValidateToken(ctx context.Context, in *sso.ValidateTokenRequest) (*sso.ValidateTokenResponse, error) {
	if in.GetAccessToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "access token is required")
	}

	isValid, username, appId, err := s.auth.ValidateToken(ctx, in.GetAccessToken())
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}
		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &sso.ValidateTokenResponse{IsValid: isValid, Username: username, AppId: int32(appId)}, nil
}
