package transport

import (
	auth "github.com/pulse-fetch/protos/pkg/pb/auth"
	"google.golang.org/grpc"
)

type Service interface {
	Register(username, email, password string) error
	Auth(username, pass string) (string, error)
	Update(id int64, newName string) error
	Del(id int64) error
}

type Auth struct {
	auth.UnimplementedAuthServer
	service Service
}

func Register(server *grpc.Server, serv Service) {
	auth.RegisterAuthServer(server, &Auth{service: serv})
}
