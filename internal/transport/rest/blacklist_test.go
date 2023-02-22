package rest

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/qnocks/blacklist-user-service/internal/entity"
	"github.com/qnocks/blacklist-user-service/internal/service"
	mock_service "github.com/qnocks/blacklist-user-service/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestHandler_saveBlacklistedUser(t *testing.T) {
	const api = "/api/blacklist"
	type mockBehavior = func(s *mock_service.MockBlacklist, user entity.BlacklistedUser)

	testTable := []struct {
		name           string
		inputBody      string
		inputUser      entity.BlacklistedUser
		mockBehavior   mockBehavior
		expectedStatus int
		expectErr      bool
	}{
		{
			name:      "OK",
			inputBody: `{"phone": "+77777777777", "username": "user", "cause": "spam", "caused_by": "other"}`,
			inputUser: entity.BlacklistedUser{
				Phone:    "+77777777777",
				Username: "user",
				Cause:    "spam",
				CausedBy: "other",
			},
			mockBehavior: func(s *mock_service.MockBlacklist, user entity.BlacklistedUser) {
				s.EXPECT().Save(user).Return(nil)
			},
			expectedStatus: 201,
			expectErr:      false,
		},
		{
			name:      "Service error",
			inputBody: `{"phone": "+77777777777", "username": "user", "cause": "spam", "caused_by": "other"}`,
			inputUser: entity.BlacklistedUser{
				Phone:    "+77777777777",
				Username: "user",
				Cause:    "spam",
				CausedBy: "other",
			},
			mockBehavior: func(s *mock_service.MockBlacklist, user entity.BlacklistedUser) {
				s.EXPECT().Save(user).Return(errors.New("error"))
			},
			expectedStatus: 500,
			expectErr:      true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockBlacklist := mock_service.NewMockBlacklist(c)
			testCase.mockBehavior(mockBlacklist, testCase.inputUser)

			services := &service.Service{Blacklist: mockBlacklist}
			handler := NewHandler(services)

			r := gin.New()
			r.POST(api, handler.saveBlacklistedUser)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, api, bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
		})
	}
}

func TestHandler_deleteBlacklistedUser(t *testing.T) {
	const api = "/api/blacklist"
	type mockBehavior = func(s *mock_service.MockBlacklist, id int)

	testTable := []struct {
		name           string
		idParam        string
		mockBehavior   mockBehavior
		expectedStatus int
		expectErr      bool
	}{
		{
			name:    "OK",
			idParam: "1",
			mockBehavior: func(s *mock_service.MockBlacklist, id int) {
				s.EXPECT().Delete(id).Return(nil)
			},
			expectedStatus: 204,
			expectErr:      false,
		},
		{
			name:    "Service error",
			idParam: "2",
			mockBehavior: func(s *mock_service.MockBlacklist, id int) {
				s.EXPECT().Delete(id).Return(errors.New("error"))
			},
			expectedStatus: 500,
			expectErr:      true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockBlacklist := mock_service.NewMockBlacklist(c)

			id, err := strconv.Atoi(testCase.idParam)
			if err != nil {
				assert.Errorf(t, err, "cannot parse id param to int - %s\n", testCase.idParam)
			}

			if testCase.mockBehavior != nil {
				testCase.mockBehavior(mockBlacklist, id)
			}

			services := &service.Service{Blacklist: mockBlacklist}
			handler := NewHandler(services)

			r := gin.New()
			r.DELETE(fmt.Sprintf("%s/:id", api), handler.deleteBlacklistedUser)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s", api, testCase.idParam), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
		})
	}
}

func TestHandler_getBlacklistedUsers(t *testing.T) {
	const api = "/api/blacklist"
	type mockBehavior = func(s *mock_service.MockBlacklist, phone, username string)

	testTable := []struct {
		name           string
		phoneParam     string
		usernameParam  string
		mockBehavior   mockBehavior
		expectedStatus int
		expectErr      bool
	}{
		{
			name:          "OK",
			phoneParam:    "89049767777",
			usernameParam: "user",
			mockBehavior: func(s *mock_service.MockBlacklist, phone, username string) {
				s.EXPECT().Find(phone, username).Return([]entity.BlacklistedUser{
					{
						ID:       1,
						Phone:    phone,
						Username: username,
						CausedBy: "John",
					},
					{
						ID:       2,
						Phone:    phone,
						Username: username,
						CausedBy: "Bob",
					},
				}, nil)
			},
			expectedStatus: 200,
			expectErr:      false,
		},
		{
			name:          "Service error",
			phoneParam:    "89049767777",
			usernameParam: "user",
			mockBehavior: func(s *mock_service.MockBlacklist, phone, username string) {
				s.EXPECT().Find(phone, username).Return([]entity.BlacklistedUser{}, errors.New("error"))
			},
			expectedStatus: 500,
			expectErr:      true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockBlacklist := mock_service.NewMockBlacklist(c)

			if testCase.mockBehavior != nil {
				testCase.mockBehavior(mockBlacklist, testCase.phoneParam, testCase.usernameParam)
			}

			services := &service.Service{Blacklist: mockBlacklist}
			handler := NewHandler(services)

			r := gin.New()
			r.GET(api, handler.getBlacklistedUsers)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s?phone=%s&username=%s",
				api, testCase.phoneParam, testCase.usernameParam), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
		})
	}
}
