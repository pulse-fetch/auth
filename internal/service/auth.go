package service

import (
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

func (s *Service) ParseJwtToken(tokenString string) (int64, error) {
	var op = "service.auth,ParseJwtToken"
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			s.log.Warn("suspicious request: invalid token signature algorithm", "actual_alg", token.Header["alg"])
			return "", fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.hmacSecret), nil
	})

	if err != nil {

		s.log.Warn("error parsing JWT token", "error", err)
		return 0, err
	}
	if !token.Valid {
		s.log.Warn("invalid JWT token")
		return 0, fmt.Errorf("%s: invalid jwt token", op)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		s.log.Warn("не удалось прочитать claims из JWT токена")
		return 0, fmt.Errorf("%s: invalid claims", op)
	}

	subject, ok := claims["sub"].(string)
	if !ok {
		s.log.Warn("invalid sub in JWT token")
		return 0, fmt.Errorf("%s: invalid sub", op)
	}

	id, err := strconv.Atoi(subject)
	if err != nil {
		slog.Warn("failed converting sub JWT in int64", "sub", subject)
		return 0, fmt.Errorf("%s: invalid converting sub", op)
	}

	return int64(id), nil
}
