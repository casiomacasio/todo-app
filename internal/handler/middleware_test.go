package handler

// import (
// 	"testing"
// 	"errors"
// 	"github.com/gin-gonic/gin"
// 	"net/http"
// 	"net/http/httptest"
// 	"github.com/casiomacasio/todo-app/internal/service"
// 	mock_service "github.com/casiomacasio/todo-app/internal/service/mock"
// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/assert"
// )

// func TestHandler_userIdentity(t *testing.T) {
// 	type mockBehavior func(s *mock_service.MockAuthorization, token string)

// 	testTable := []struct {
// 		name               string
// 		header             string
// 		token              string
// 		mockBehavior       mockBehavior
// 		expectedStatusCode int
// 		expectedBody       string
// 	}{
// 		{
// 			name:               "No header",
// 			header:             "",
// 			mockBehavior:       func(s *mock_service.MockAuthorization, token string) {},
// 			expectedStatusCode: http.StatusUnauthorized,
// 			expectedBody:       `{"message":"empty auth header"}`,
// 		},
// 		{
// 			name:               "Invalid header format",
// 			header:             "BearerTokenOnly",
// 			mockBehavior:       func(s *mock_service.MockAuthorization, token string) {},
// 			expectedStatusCode: http.StatusUnauthorized,
// 			expectedBody:       `{"message":"header's length in invalid"}`,
// 		},
// 		{
// 			name:   "Invalid token",
// 			header: "Bearer invalidtoken",
// 			token:  "invalidtoken",
// 			mockBehavior: func(s *mock_service.MockAuthorization, token string) {
// 				s.EXPECT().ParseToken(token).Return(0, errors.New("invalid token"))
// 			},
// 			expectedStatusCode: http.StatusUnauthorized,
// 			expectedBody:       `{"message":"invalid token"}`,
// 		},
// 		{
// 			name:   "Success",
// 			header: "Bearer validtoken",
// 			token:  "validtoken",
// 			mockBehavior: func(s *mock_service.MockAuthorization, token string) {
// 				s.EXPECT().ParseToken(token).Return(42, nil)
// 			},
// 			expectedStatusCode: http.StatusOK,
// 			expectedBody:       `{"message":"ok"}`,
// 		},
// 	}

// 	for _, testCase := range testTable {
// 		t.Run(testCase.name, func(t *testing.T) {
// 			gin.SetMode(gin.TestMode)
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			auth := mock_service.NewMockAuthorization(ctrl)
// 			testCase.mockBehavior(auth, testCase.token)

// 			services := &service.Service{Authorization: auth}
// 			handler := NewHandler(services)
// 			router := gin.New()
// 			router.GET("/protected", handler.userIdentity, func(c *gin.Context) {
// 				c.JSON(http.StatusOK, gin.H{"message": "ok"})
// 			})

// 			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
// 			if testCase.header != "" {
// 				req.Header.Set("Authorization", testCase.header)
// 			}
// 			w := httptest.NewRecorder()
// 			router.ServeHTTP(w, req)

// 			assert.Equal(t, testCase.expectedStatusCode, w.Code)
// 			assert.JSONEq(t, testCase.expectedBody, w.Body.String())
// 		})
// 	}
// }

// func Test_getUserID(t *testing.T) {
// 	testTable := []struct {
// 		name          string
// 		ctxSetup      func(*gin.Context)
// 		expectedID    int
// 		expectedError string
// 	}{
// 		{
// 			name: "Success",
// 			ctxSetup: func(c *gin.Context) {
// 				c.Set(userCtx, 123)
// 			},
// 			expectedID:    123,
// 			expectedError: "",
// 		},
// 		{
// 			name: "Not found",
// 			ctxSetup: func(c *gin.Context) {},
// 			expectedID:    0,
// 			expectedError: "user id not found",
// 		},
// 		{
// 			name: "Wrong type",
// 			ctxSetup: func(c *gin.Context) {
// 				c.Set(userCtx, "not an int")
// 			},
// 			expectedID:    0,
// 			expectedError: "user id is not of valid type",
// 		},
// 	}

// 	for _, tt := range testTable {
// 		t.Run(tt.name, func(t *testing.T) {
// 			gin.SetMode(gin.TestMode)
// 			c, _ := gin.CreateTestContext(httptest.NewRecorder())
// 			tt.ctxSetup(c)

// 			id, err := getUserID(c)
// 			assert.Equal(t, tt.expectedID, id)
// 			if tt.expectedError == "" {
// 				assert.NoError(t, err)
// 			} else {
// 				assert.EqualError(t, err, tt.expectedError)
// 			}
// 		})
// 	}
// }
