package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/qnocks/blacklist-user-service/internal/entity"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type Blacklist interface {
	Save(user entity.BlacklistedUser) error
	Delete(id int) error
	Find(phone, username string) ([]entity.BlacklistedUser, error)
}

type Auth interface {
	GetUser(username, password string) (entity.User, error)
}

type Repository struct {
	Blacklist
	Auth
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Blacklist: NewBlacklistRepository(db),
		Auth:      NewAuthRepository(db),
	}
}
