package service

import (
	"auth/internal/domain"
	hash "auth/pkg/hasher"
	"log/slog"
)

type Storage interface {
	Create(username, email, hashPass string) error
	Get(data interface{}, method string) (domain.User, error)
	Del(id int64) error
	UpdateName(id int64, newName string) error
}

type Service struct {
	log        *slog.Logger
	storage    Storage
	hasher     *hash.Sha256
	hmacSecret string
}

func New(log *slog.Logger, storage Storage, hasher *hash.Sha256, hmacSecret string) *Service {
	return &Service{log: log, storage: storage, hasher: hasher, hmacSecret: hmacSecret}
}
