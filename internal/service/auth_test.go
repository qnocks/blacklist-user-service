package service

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/qnocks/blacklist-user-service/internal/entity"
	mock_repository "github.com/qnocks/blacklist-user-service/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthService_Login(t *testing.T) {
	type mockBehavior = func(s *mock_repository.MockAuth, user entity.User)

	testTable := []struct {
		name         string
		user         entity.User
		token        string
		mockBehavior mockBehavior
		expectErr    bool
		errorMessage string
	}{
		{
			name: "OK",
			user: entity.User{
				Username: "user",
				Password: "user",
			},
			mockBehavior: func(s *mock_repository.MockAuth, user entity.User) {
				s.EXPECT().GetUser(user.Username, user.Password).Return(user, nil)
			},
			expectErr: false,
		},
		{
			name: "Repository error",
			user: entity.User{
				Username: "user",
				Password: "user",
			},
			mockBehavior: func(s *mock_repository.MockAuth, user entity.User) {
				s.EXPECT().GetUser(user.Username, user.Password).Return(entity.User{}, errors.New("error"))
			},
			expectErr:    true,
			errorMessage: "error",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)

			mockAuth := mock_repository.NewMockAuth(c)
			authService := NewAuthService(mockAuth)

			testCase.mockBehavior(mockAuth, entity.User{
				Username: testCase.user.Username,
				Password: authService.generatePasswordHash(testCase.user.Password),
			})

			_, err := authService.Login(testCase.user)

			if testCase.expectErr {
				assert.Equal(t, err.Error(), testCase.errorMessage)
			} else {
				assert.Equal(t, err, nil)
			}
		})
	}
}
