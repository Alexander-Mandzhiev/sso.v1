package app

import (
	grpcapp "sso/internal/app/grpc"
	"sso/internal/repository"
	"sso/internal/service"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(upp, agroReports, session *repository.Repository) *App {
	authService := service.New(upp, agroReports, session)
	grpcApp := grpcapp.New(authService)
	return &App{
		GRPCServer: grpcApp,
	}
}
