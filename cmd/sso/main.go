package main

import (
	"log/slog"
	"os"
	"os/signal"
	"sso/internal/app"
	"sso/internal/config"
	"sso/internal/repository"
	sl "sso/pkg/logger"
	"syscall"
)

func main() {
	config.Initialize()
	sl.SetupLogger()
	sl.Log.Info("Starting service sso", slog.String("address", config.Cfg.Address), slog.Int("port", config.Cfg.Port))
	sl.Log.Debug("Debug messages are enabled")

	upp, err := repository.New(config.Cfg.DBUpp)
	if err != nil {
		panic(err)
	}
	agroReports, err := repository.New(config.Cfg.DBAgroReports)
	if err != nil {
		panic(err)
	}

	session, err := repository.New(config.Cfg.DBAgroReports)

	application := app.New(upp, agroReports, session)
	go func() {
		application.GRPCServer.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCServer.Stop()

	if err = upp.Stop(); err != nil {
		sl.Log.Error("Failed to stop UPP repository", slog.Any("error", err))
	}
	if err = agroReports.Stop(); err != nil {
		sl.Log.Error("Failed to stop AgroReports repository", slog.Any("error", err))
	}

	sl.Log.Info("Gracefully stopped")
}
