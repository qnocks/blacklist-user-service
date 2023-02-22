package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/qnocks/blacklist-user-service/internal/entity"
	"time"
)

const (
	blacklistedUsersTable                   = "blacklisted_users"
	selectBlacklistedUsersByPhoneOrUsername = "SELECT * FROM " + blacklistedUsersTable + " WHERE phone=$1 OR username=$2"
	deleteBlacklistedUserByID               = "DELETE FROM " + blacklistedUsersTable + " WHERE id=$1"
	insertBlacklistedUser                   = "INSERT INTO " + blacklistedUsersTable +
		" (phone, username, cause, timestamp, caused_by) VALUES ($1, $2, $3, $4, $5)"
)

type BlacklistRepository struct {
	DB *sqlx.DB
}

func (r *BlacklistRepository) Save(user entity.BlacklistedUser) error {
	user.Timestamp = time.Now()
	row := r.DB.QueryRow(insertBlacklistedUser, user.Phone, user.Username, user.Cause, user.Timestamp, user.CausedBy)
	return row.Err()
}

func (r *BlacklistRepository) Delete(id int) error {
	row := r.DB.QueryRow(deleteBlacklistedUserByID, id)
	return row.Err()
}

func (r *BlacklistRepository) Find(phone, username string) ([]entity.BlacklistedUser, error) {
	var users []entity.BlacklistedUser
	err := r.DB.Select(&users, selectBlacklistedUsersByPhoneOrUsername, phone, username)
	return users, err
}

func NewBlacklistRepository(DB *sqlx.DB) *BlacklistRepository {
	return &BlacklistRepository{DB: DB}
}
