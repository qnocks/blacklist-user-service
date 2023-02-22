package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/qnocks/blacklist-user-service/internal/entity"
	"github.com/qnocks/blacklist-user-service/internal/repository"
	"os"
	"strconv"
	"time"
)

type AuthService struct {
	repo        repository.Auth
	salt        string
	tokenSecret string
	tokenTTL    time.Duration
}

func NewAuthService(repo repository.Auth) *AuthService {
	ttl := os.Getenv("TOKEN_TTL")
	tokenTTL, _ := strconv.Atoi(ttl)
	return &AuthService{
		repo:        repo,
		salt:        os.Getenv("PASSWORD_SALT"),
		tokenSecret: os.Getenv("TOKEN_SECRET"),
		tokenTTL:    time.Duration(tokenTTL),
	}
}

func (s *AuthService) Login(user entity.User) (string, error) {
	user, err := s.repo.GetUser(user.Username, s.generatePasswordHash(user.Password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(s.tokenTTL * time.Minute).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   user.Username,
	})

	return token.SignedString([]byte(s.tokenSecret))
}

func (s *AuthService) ParseToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(s.tokenSecret), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok && !token.Valid {
		return "", err
	}

	return claims.Subject, nil
}

func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(s.salt)))
}
