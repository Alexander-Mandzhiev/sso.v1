package auth

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"sso/internal/service"
	"sso/protos/gen/go/sso"
)

func (s *serverAPI) SignUp(ctx context.Context, in *sso.SignupRequest) (*sso.SignupResponse, error) {
	if in.GetUsername() == "" || in.GetPassword() == "" || in.GetAppId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "username, password, and app_id are required")
	}

	accessToken, refreshToken, err := s.auth.SignUp(ctx, in.GetUsername(), in.GetPassword(), int(in.GetAppId()))
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}

		return nil, status.Error(codes.Internal, "failed to login")
	}

	outgoingMD := metadata.Pairs("refresh_token", refreshToken)
	grpc.SendHeader(ctx, outgoingMD)

	return &sso.SignupResponse{AccessToken: accessToken}, nil
}
