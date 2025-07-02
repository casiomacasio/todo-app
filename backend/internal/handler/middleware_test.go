package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"github.com/redis/go-redis/v9"
	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRedis() (*redis.Client, func()) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	return client, func() {
		client.Close()
		s.Close()
	}
}

func setupRouterWithMiddleware(middleware gin.HandlerFunc, userCtxEnabled bool) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		if userCtxEnabled {
			c.Set("userId", 1)
		}
		c.Next()
	})

	r.Use(middleware)
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	return r
}

func TestRateLimitMiddleware(t *testing.T) {
	redisClient, cleanup := setupTestRedis()
	defer cleanup()

	router := setupRouterWithMiddleware(RateLimitMiddleware(redisClient, 3, time.Minute), true)

	for i := 1; i <= 4; i++ {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		router.ServeHTTP(rec, req)

		if i <= 3 {
			assert.Equal(t, http.StatusOK, rec.Code, "request #%d should pass", i)
		} else {
			assert.Equal(t, http.StatusTooManyRequests, rec.Code, "request #%d should be rate limited", i)
		}
	}
}

func TestGlobalRateLimitMiddleware(t *testing.T) {
	redisClient, cleanup := setupTestRedis()
	defer cleanup()

	router := setupRouterWithMiddleware(GlobalRateLimitMiddleware(redisClient, 2, time.Minute), false)

	for i := 1; i <= 3; i++ {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		router.ServeHTTP(rec, req)

		if i <= 2 {
			assert.Equal(t, http.StatusOK, rec.Code, "request #%d should pass", i)
		} else {
			assert.Equal(t, http.StatusTooManyRequests, rec.Code, "request #%d should be rate limited", i)
		}
	}
}

func TestRateLimitIPMiddleware(t *testing.T) {
	redisClient, cleanup := setupTestRedis()
	defer cleanup()

	router := setupRouterWithMiddleware(RateLimitIPMiddleware(redisClient, 1, time.Minute), false)

	rec1 := httptest.NewRecorder()
	req1, _ := http.NewRequest(http.MethodGet, "/", nil)
	req1.RemoteAddr = "1.2.3.4:1234"
	router.ServeHTTP(rec1, req1)
	assert.Equal(t, http.StatusOK, rec1.Code)

	rec2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodGet, "/", nil)
	req2.RemoteAddr = "1.2.3.4:1234"
	router.ServeHTTP(rec2, req2)
	assert.Equal(t, http.StatusTooManyRequests, rec2.Code)
}
