package rest

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/qnocks/blacklist-user-service/internal/entity"
	"github.com/qnocks/blacklist-user-service/internal/service"
	mock_service "github.com/qnocks/blacklist-user-service/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_login(t *testing.T) {
	const api = "/auth/login"
	type mockBehavior = func(s *mock_service.MockAuth, user entity.User)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            entity.User
		mockBehavior         mockBehavior
		expectedStatus       int
		expectErr            bool
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"username": "test", "password": "test"}`,
			inputUser: entity.User{
				Username: "test",
				Password: "test",
			},
			mockBehavior: func(s *mock_service.MockAuth, user entity.User) {
				s.EXPECT().Login(user).Return("generated-token", nil)
			},
			expectedStatus:       200,
			expectedResponseBody: `{"token":"generated-token"}`,
			expectErr:            false,
		},
		{
			name:      "Service error",
			inputBody: `{"username": "test", "password": "incorrect"}`,
			inputUser: entity.User{
				Username: "test",
				Password: "incorrect",
			},
			mockBehavior: func(s *mock_service.MockAuth, user entity.User) {
				s.EXPECT().Login(user).Return("", errors.New("error"))
			},
			expectedStatus: 400,
			expectErr:      true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockAuth := mock_service.NewMockAuth(c)
			testCase.mockBehavior(mockAuth, testCase.inputUser)

			services := &service.Service{Auth: mockAuth}
			handler := NewHandler(services)

			r := gin.New()
			r.POST(api, handler.login)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, api, bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
			if !testCase.expectErr {
				assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
			}
		})
	}
}
