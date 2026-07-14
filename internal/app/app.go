package app

import (
	grpcapp "auth/internal/app/grpc"
	"auth/internal/domain"
	"auth/internal/service"
	"auth/internal/storage"
	"auth/internal/storage/postgres"
	"log/slog"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, port string, params domain.PostgresParams) *App {
	db, err := postgres.Conn(params)
	if err != nil {
		panic("Failed starting postgres, error: " + err.Error())
	}
	repo := storage.New(db)
	serv := service.New(log, repo)
	server := grpcapp.New(log, port, serv)
	return &App{GRPCServer: server}
}
