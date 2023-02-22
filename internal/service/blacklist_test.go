package service

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/qnocks/blacklist-user-service/internal/entity"
	mock_repository "github.com/qnocks/blacklist-user-service/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBlacklistService_Save(t *testing.T) {
	type mockBehavior = func(s *mock_repository.MockBlacklist, user entity.BlacklistedUser)

	testTable := []struct {
		name         string
		user         entity.BlacklistedUser
		mockBehavior mockBehavior
		expectErr    bool
		errorMessage string
	}{
		{
			name: "OK",
			user: entity.BlacklistedUser{
				Phone:    "+77777777777",
				Username: "user",
				Cause:    "spam",
				CausedBy: "other",
			},
			mockBehavior: func(s *mock_repository.MockBlacklist, user entity.BlacklistedUser) {
				s.EXPECT().Save(user).Return(nil)
			},
			expectErr: false,
		},
		{
			name: "Repository error",
			user: entity.BlacklistedUser{
				Phone:    "+77777777777",
				Username: "user",
				Cause:    "spam",
				CausedBy: "other",
			},
			mockBehavior: func(s *mock_repository.MockBlacklist, user entity.BlacklistedUser) {
				s.EXPECT().Save(user).Return(errors.New("error"))
			},
			expectErr:    true,
			errorMessage: "error",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)

			mockLink := mock_repository.NewMockBlacklist(c)
			testCase.mockBehavior(mockLink, testCase.user)

			blacklistService := NewBlacklistService(mockLink)
			err := blacklistService.Save(testCase.user)

			if testCase.expectErr {
				assert.Equal(t, err.Error(), testCase.errorMessage)
			} else {
				assert.Equal(t, err, nil)
			}
		})
	}
}

func TestBlacklistService_Delete(t *testing.T) {
	type mockBehavior = func(s *mock_repository.MockBlacklist, id int)

	testTable := []struct {
		name         string
		id           int
		mockBehavior mockBehavior
		expectErr    bool
		errorMessage string
	}{
		{
			name: "OK",
			id:   1,
			mockBehavior: func(s *mock_repository.MockBlacklist, id int) {
				s.EXPECT().Delete(id).Return(nil)
			},
			expectErr: false,
		},
		{
			name: "Repository error",
			id:   1,
			mockBehavior: func(s *mock_repository.MockBlacklist, id int) {
				s.EXPECT().Delete(id).Return(errors.New("error"))
			},
			expectErr:    true,
			errorMessage: "error",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)

			mockLink := mock_repository.NewMockBlacklist(c)
			testCase.mockBehavior(mockLink, testCase.id)

			blacklistService := NewBlacklistService(mockLink)
			err := blacklistService.Delete(testCase.id)

			if testCase.expectErr {
				assert.Equal(t, err.Error(), testCase.errorMessage)
			} else {
				assert.Equal(t, err, nil)
			}
		})
	}
}

func TestBlacklistService_Find(t *testing.T) {
	type mockBehavior = func(s *mock_repository.MockBlacklist, phone, username string)

	testTable := []struct {
		name         string
		phone        string
		username     string
		users        []entity.BlacklistedUser
		mockBehavior mockBehavior
		expectErr    bool
		errorMessage string
	}{
		{
			name:  "OK",
			phone: "+77777777777",
			users: []entity.BlacklistedUser{
				{
					ID:       1,
					Phone:    "+77777777777",
					Username: "user",
					Cause:    "spam",
					CausedBy: "other",
				},
				{
					ID:       2,
					Phone:    "+77777777777",
					Username: "another",
					Cause:    "bullying",
					CausedBy: "other",
				},
			},
			mockBehavior: func(s *mock_repository.MockBlacklist, phone, username string) {
				s.EXPECT().Find(phone, username).Return([]entity.BlacklistedUser{
					{
						ID:       1,
						Phone:    "+77777777777",
						Username: "user",
						Cause:    "spam",
						CausedBy: "other",
					},
					{
						ID:       2,
						Phone:    "+77777777777",
						Username: "another",
						Cause:    "bullying",
						CausedBy: "other",
					},
				}, nil)
			},
			expectErr: false,
		},
		{
			name:  "Repository error",
			phone: "+77777777777",
			users: []entity.BlacklistedUser{},
			mockBehavior: func(s *mock_repository.MockBlacklist, phone, username string) {
				s.EXPECT().Find(phone, username).Return([]entity.BlacklistedUser{}, errors.New("error"))
			},
			expectErr:    true,
			errorMessage: "error",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)

			mockLink := mock_repository.NewMockBlacklist(c)
			testCase.mockBehavior(mockLink, testCase.phone, testCase.username)

			blacklistService := NewBlacklistService(mockLink)
			actual, err := blacklistService.Find(testCase.phone, testCase.username)

			if testCase.expectErr {
				assert.Equal(t, err.Error(), testCase.errorMessage)
			} else {
				assert.Equal(t, err, nil)
				assert.Equal(t, actual, testCase.users)
			}
		})
	}

}
