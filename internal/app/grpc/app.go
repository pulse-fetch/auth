package grpcapp

import (
	"auth/internal/service"
	"auth/internal/transport"
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type App struct {
	log        *slog.Logger
	GRPCServer *grpc.Server
	port       string
}

func New(log *slog.Logger, port string, serv *service.Service) *App {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.PayloadReceived, logging.PayloadSent,
		),
	}
	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			log.Error("Recovered from panic", slog.Any("panic", p))

			return status.Errorf(codes.Internal, "internal error")
		}),
	}
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		logging.UnaryServerInterceptor(InterceptorLogger(log), loggingOpts...),
	))
	transport.Register(grpcServer, serv)

	return &App{GRPCServer: grpcServer, port: port, log: log}
}

func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic("Failed starting server, error: " + err.Error())
	}
}

func (a *App) Run() error {
	var op = "grpcapp.Run"
	addr := fmt.Sprintf(":%s", a.port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}
	a.log.Info("Server starting...", slog.String("op", op), slog.String("port", a.port))
	if err = a.GRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}
	return nil

}

func (a *App) Stop() {
	var op = "grpcapp.Stop"
	a.log.Info("Server stopping...", slog.String("op", op))
	a.GRPCServer.GracefulStop()
}
