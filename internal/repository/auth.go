package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/qnocks/blacklist-user-service/internal/entity"
)

const (
	authTable                       = "users"
	selectUserByUsernameAndPassword = "SELECT * FROM " + authTable + " WHERE username=$1 AND password=$2"
)

type AuthRepository struct {
	DB *sqlx.DB
}

func (r *AuthRepository) GetUser(username, password string) (entity.User, error) {
	var user entity.User
	err := r.DB.Get(&user, selectUserByUsernameAndPassword, username, password)
	return user, err
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{DB: db}
}
