package service

import (
	"github.com/qnocks/blacklist-user-service/internal/entity"
	"github.com/qnocks/blacklist-user-service/internal/repository"
)

type BlacklistService struct {
	repo repository.Blacklist
}

func NewBlacklistService(repo repository.Blacklist) *BlacklistService {
	return &BlacklistService{repo: repo}
}

func (s *BlacklistService) Save(user entity.BlacklistedUser) error {
	return s.repo.Save(user)
}

func (s *BlacklistService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *BlacklistService) Find(phone, username string) ([]entity.BlacklistedUser, error) {
	return s.repo.Find(phone, username)
}
