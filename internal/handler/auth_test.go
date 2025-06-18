package handler

// import (
// 	"testing"
// 	"errors"
// 	"github.com/gin-gonic/gin"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"github.com/casiomacasio/todo-app/internal/domain"
// 	"github.com/casiomacasio/todo-app/internal/repository"
// 	"github.com/casiomacasio/todo-app/internal/service"
// 	mock_service "github.com/casiomacasio/todo-app/internal/service/mock"
// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/assert"
// )


// func TestHandlerSignUp(t *testing.T) {
// 	type mockBehavior func(s *mock_service.MockAuthorization, user domain.User)

// 	gin.SetMode(gin.TestMode)

// 	testTable := []struct {
// 		name                string
// 		inputBody           string
// 		inputUser           domain.User
// 		mockBehavior        mockBehavior
// 		expectedStatusCode  int
// 		expectedRequestBody string
// 	}{
// 		{
// 			name:      "ok",
// 			inputBody: `{"name":"string","username":"string","password":"string"}`,
// 			inputUser: domain.User{
// 				Name:     "string",
// 				Username: "string",
// 				Password: "string",
// 			},
// 			mockBehavior: func(s *mock_service.MockAuthorization, user domain.User) {
// 				s.EXPECT().CreateUser(user).Return(1, nil)
// 			},
// 			expectedStatusCode:  200,
// 			expectedRequestBody: "{\"id\":1}",
// 		},
// 		{
// 			name:                "empty fields",
// 			inputBody:           `{}`,
// 			inputUser:           domain.User{},
// 			mockBehavior:        func(s *mock_service.MockAuthorization, user domain.User) {},
// 			expectedStatusCode:  400,
// 			expectedRequestBody: "{\"message\":\"invalid body\"}",
// 		},
// 		{
// 			name:      "username already in use",
// 			inputBody: `{"name":"string","username":"string","password":"string"}`,
// 			inputUser: domain.User{
// 				Name:     "string",
// 				Username: "string",
// 				Password: "string",
// 			},
// 			mockBehavior: func(s *mock_service.MockAuthorization, user domain.User) {
// 				s.EXPECT().CreateUser(user).Return(0, repository.ErrUsernameExists)
// 			},
// 			expectedStatusCode:  409,
// 			expectedRequestBody: "{\"message\":\"Username already in use\"}",
// 		},
// 	}

// 	for _, testCase := range testTable {
// 		t.Run(testCase.name, func(t *testing.T) {
// 			c := gomock.NewController(t)
// 			defer c.Finish()

// 			auth := mock_service.NewMockAuthorization(c)
// 			testCase.mockBehavior(auth, testCase.inputUser)

// 			services := &service.Service{Authorization: auth}
// 			handler := NewHandler(services)

// 			e := gin.New()
// 			e.POST("/sign-up", handler.signUp)

// 			req := httptest.NewRequest(http.MethodPost, "/sign-up", strings.NewReader(testCase.inputBody))
// 			req.Header.Set("Content-Type", "application/json") 
// 			rec := httptest.NewRecorder()

// 			e.ServeHTTP(rec, req)

// 			assert.Equal(t, testCase.expectedStatusCode, rec.Code)
// 			assert.Equal(t, testCase.expectedRequestBody, rec.Body.String())
// 		})
// 	}
// }

// func TestHandlerSignIn(t *testing.T) {
// 	type mockBehavior func(s *mock_service.MockAuthorization, username, password string)

// 	gin.SetMode(gin.TestMode)

// 	testTable := []struct {
// 		name                string
// 		inputBody           string
// 		username            string
// 		password            string
// 		mockBehavior        mockBehavior
// 		expectedStatusCode  int
// 		expectedErrorSubstr []string
// 		expectedResponse    string
// 	}{
// 		{
// 			name:     "ok",
// 			inputBody: `{"username":"user1","password":"pass1"}`,
// 			username:  "user1",
// 			password:  "pass1",
// 			mockBehavior: func(s *mock_service.MockAuthorization, username, password string) {
// 				s.EXPECT().GenerateToken(username, password).Return("valid-token", nil)
// 			},
// 			expectedStatusCode: 200,
// 			expectedResponse:   `{"token":"valid-token"}`,
// 		},
// 		{
// 			name:     "empty fields",
// 			inputBody: `{}`,
// 			mockBehavior: func(s *mock_service.MockAuthorization, username, password string) {},
// 			expectedStatusCode: 400,
// 			expectedErrorSubstr: []string{
// 				"Username", "Password", "required",
// 			},
// 		},
// 		{
// 			name:     "invalid credentials",
// 			inputBody: `{"username":"wrong","password":"bad"}`,
// 			username:  "wrong",
// 			password:  "bad",
// 			mockBehavior: func(s *mock_service.MockAuthorization, username, password string) {
// 				s.EXPECT().GenerateToken(username, password).Return("", errors.New("invalid credentials"))
// 			},
// 			expectedStatusCode: 500,
// 			expectedResponse:   `{"message":"invalid credentials"}`,
// 		},
// 	}

// 	for _, testCase := range testTable {
// 		t.Run(testCase.name, func(t *testing.T) {
// 			c := gomock.NewController(t)
// 			defer c.Finish()

// 			auth := mock_service.NewMockAuthorization(c)
// 			testCase.mockBehavior(auth, testCase.username, testCase.password)

// 			services := &service.Service{Authorization: auth}
// 			handler := NewHandler(services)

// 			e := gin.New()
// 			e.POST("/sign-in", handler.signIn)

// 			req := httptest.NewRequest(http.MethodPost, "/sign-in", strings.NewReader(testCase.inputBody))
// 			req.Header.Set("Content-Type", "application/json")
// 			rec := httptest.NewRecorder()

// 			e.ServeHTTP(rec, req)

// 			assert.Equal(t, testCase.expectedStatusCode, rec.Code)

// 			if len(testCase.expectedErrorSubstr) > 0 {
// 				for _, substr := range testCase.expectedErrorSubstr {
// 					assert.Contains(t, rec.Body.String(), substr)
// 				}
// 			} else {
// 				assert.JSONEq(t, testCase.expectedResponse, rec.Body.String())
// 			}
// 		})
// 	}
// }
