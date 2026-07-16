package transport

import (
	"context"

	auth "github.com/pulse-fetch/protos/pkg/pb/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *Auth) Register(ctx context.Context, r *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	err := a.service.Register(r.GetUsername(), r.GetEmail(), r.GetPassword())
	if err != nil {
		return &auth.RegisterResponse{}, err
	}
	return &auth.RegisterResponse{Status: "OK"}, nil
}

func (a *Auth) Login(ctx context.Context, r *auth.LoginRequest) (*auth.LoginResponse, error) {

	jwt, err := a.service.Auth(r.GetUsername(), r.GetPassword())
	if err != nil {
		return &auth.LoginResponse{}, status.Error(codes.Internal, "failed to login")
	}
	return &auth.LoginResponse{Jwt: jwt}, nil
}

func (a *Auth) UpdateName(ctx context.Context, r *auth.UpdateRequest) (*auth.UpdateResponse, error) {

	if err := a.service.Update(r.GetId(), r.GetNewName()); err != nil {

		return &auth.UpdateResponse{Status: "false"}, status.Error(codes.Internal, "failed update name")
	}
	return &auth.UpdateResponse{Status: "OK"}, nil
}

func (a *Auth) Delete(ctx context.Context, r *auth.DeleteRequest) (*auth.DeleteResponse, error) {
	if err := a.service.Del(r.GetId()); err != nil {
		return &auth.DeleteResponse{Status: "false"}, status.Error(codes.Internal, "failed delete user")
	}
	return &auth.DeleteResponse{Status: "OK"}, nil
}
func (a *Auth) Get(ctx context.Context, r *auth.GetRequest) (*auth.GetResponse, error) {
	resp, err := a.service.Get(r.GetId())
	if err != nil {
		return &auth.GetResponse{}, status.Error(codes.Internal, "failed delete user")
	}
	return &auth.GetResponse{Username: resp.Username, Email: resp.Email}, nil
}
