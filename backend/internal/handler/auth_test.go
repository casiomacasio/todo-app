package handler

import (
	"net/http"

	"net/http/httptest"
	"strings"
	"testing"
	"github.com/redis/go-redis/v9"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
	"github.com/alicebob/miniredis/v2" 
	"github.com/casiomacasio/todo-app/backend/internal/domain"
	"github.com/casiomacasio/todo-app/backend/internal/repository"
	"github.com/casiomacasio/todo-app/backend/internal/service"
	mock_service "github.com/casiomacasio/todo-app/backend/internal/service/mocks"
)

func TestHandlerSignUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user domain.CreateUserRequest)

	gin.SetMode(gin.TestMode)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           domain.CreateUserRequest
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "ok",
			inputBody: `{"name":"string","username":"string","password":"string"}`,
			inputUser: domain.CreateUserRequest{
				Name:     "string",
				Username: "string",
				Password: "string",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user domain.CreateUserRequest) {
				s.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "{\"id\":1}",
		},
		{
			name:                "empty fields",
			inputBody:           `{}`,
			inputUser:           domain.CreateUserRequest{},
			mockBehavior:        func(s *mock_service.MockAuthorization, user domain.CreateUserRequest) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"message\":\"invalid body\"}",
		},
		{
			name:      "username already in use",
			inputBody: `{"name":"string","username":"string","password":"string"}`,
			inputUser: domain.CreateUserRequest{
				Name:     "string",
				Username: "string",
				Password: "string",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user domain.CreateUserRequest) {
				s.EXPECT().CreateUser(user).Return(0, repository.ErrUsernameExists)
			},
			expectedStatusCode:  409,
			expectedRequestBody: "{\"message\":\"Username already in use\"}",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			srv, err := miniredis.Run()
			if err != nil {
				t.Fatalf("Failed to start miniredis: %v", err)
			}
			rdb := redis.NewClient(&redis.Options{
				Addr: srv.Addr(),
			})
			defer srv.Close()
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services, rdb)

			e := gin.New()
			e.POST("/sign-up", handler.signUp)

			req := httptest.NewRequest(http.MethodPost, "/sign-up", strings.NewReader(testCase.inputBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, testCase.expectedStatusCode, rec.Code)
			assert.Equal(t, testCase.expectedRequestBody, rec.Body.String())
		})
	}
}
