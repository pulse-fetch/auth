package transport

import (
	"context"
	"fmt"

	auth "github.com/pulse-fetch/protos/pkg/pb/auth"
)

func (a *Auth) Register(ctx context.Context, r *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	resp, err := a.service.Register(r.GetUsername(), r.GetEmail(), r.GetPassword())
	fmt.Println(resp, err)
	return &auth.RegisterResponse{}, nil
}
