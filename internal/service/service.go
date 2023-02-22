package service

import (
	"github.com/qnocks/blacklist-user-service/internal/entity"
	"github.com/qnocks/blacklist-user-service/internal/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Auth interface {
	Login(user entity.User) (string, error)
	ParseToken(token string) (string, error)
}

type Blacklist interface {
	Save(user entity.BlacklistedUser) error
	Delete(id int) error
	Find(phone, username string) ([]entity.BlacklistedUser, error)
}

type Service struct {
	Blacklist
	Auth
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Auth:      NewAuthService(repos.Auth),
		Blacklist: NewBlacklistService(repos.Blacklist),
	}
}
