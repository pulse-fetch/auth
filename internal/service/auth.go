package service

import (
	"auth/internal/domain"
	hash "auth/pkg/hasher"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (s *Service) Register(username, email, password string) error {
	var op = "service.auth.Register"
	hashPassword := hash.Hash(*s.hasher, password)
	if err := s.storage.Create(username, email, hashPassword); err != nil {
		s.log.Error("Failed create new user", slog.Any("error", err))
		return fmt.Errorf("%s:%v", op, err)
	}
	return nil
}
func (s *Service) Del(id int64) error {
	var op = "service.Auth.Del"
	if err := s.storage.Del(id); err != nil {
		s.log.Error("Failed update name", slog.Any("error", err), slog.String("op", op))
		return fmt.Errorf("%s:%v", op, err)
	}
	return nil
}

func (s *Service) Update(id int64, newName string) error {
	var op = "service.Auth.Update"
	if err := s.storage.UpdateName(id, newName); err != nil {
		s.log.Error("Failed update name", slog.Any("error", err), slog.String("op", op))
		return fmt.Errorf("%s:%v", op, err)
	}
	return nil

}

func (s *Service) Get(id int64) (domain.ResponseUser, error) {
	var op = "service.auth.Get"
	user, err := s.storage.Get(id, "id")
	if err != nil {
		s.log.Warn("Failed get user", slog.Any("error", err), slog.String("op", op))
		return domain.ResponseUser{}, err
	}
	resp := domain.ResponseUser{Username: user.Username, Email: user.Email}
	return resp, nil
}
func (s *Service) Auth(username, pass string) (string, error) {
	var op = "service.Auth"
	hashPass := hash.Hash(*s.hasher, pass)
	user, err := s.storage.Get(username, "name")
	if err != nil {
		return "", fmt.Errorf("%s: %v", op, err)
	}

	if hashPass != user.HashPassword {
		s.log.Warn("Unsuccessful attempt: incorrect password.", "username", username, "user_id", user.Id)
		return "", fmt.Errorf("%s: incorret password", op)
	}

	tokenString, err := s.generateJwtToken(user.Id)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *Service) generateJwtToken(userID int64) (string, error) {
	var op = "service.auth.generateJwtToken"
	mapClaims := jwt.MapClaims{
		"sub": strconv.FormatInt(userID, 10),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	tokenString, err := token.SignedString([]byte(s.hmacSecret))
	if err != nil {
		s.log.Error("critical error generation jwt token", "user_id", userID, "error", err)
		return "", fmt.Errorf("%s:%v", op, err)
	}
	return tokenString, nil
}
