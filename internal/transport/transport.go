package transport

import (
	auth "github.com/pulse-fetch/protos/pkg/pb/auth"
	"google.golang.org/grpc"
)

type Service interface {
	Register(username, email, password string) (string, error)
}

type Auth struct {
	auth.UnimplementedAuthServer
	service Service
}

func Register(server *grpc.Server, serv Service) {
	auth.RegisterAuthServer(server, &Auth{service: serv})
}
