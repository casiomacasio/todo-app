package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/casiomacasio/todo-app/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const userCtx = "userId"

var ctx = context.Background()

func (h *Handler) userIdentity(c *gin.Context) {
	token, err := c.Cookie("access_token")
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "missing access token cookie")
		return
	}

	userID, err := h.service.Authorization.ParseToken(token)
	if err != nil {
		if errors.Is(err, service.ErrTokenExpired) {
			c.Header("Token-Expired", "true")
			newErrorResponse(c, http.StatusUnauthorized, "token expired")
			return
		}
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(userCtx, userID)
}

func getUserID(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return 0, errors.New("user id not found")
	}
	userID, ok := id.(int)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id invalid")
		return 0, errors.New("user id invalid")
	}
	return userID, nil
}

func RateLimitMiddleware(redisClient *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := getUserID(c)
		if err != nil {
			c.Abort()
			return
		}

		key := fmt.Sprintf("rate_limit:%d", userID)
		count, err := redisClient.Incr(ctx, key).Result()
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, "redis error")
			c.Abort()
			return
		}

		if count == 1 {
			redisClient.Expire(ctx, key, window)
		}

		if count > int64(limit) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}

		c.Next()
	}
}

func GlobalRateLimitMiddleware(redisClient *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := "rate_limit:global"
		count, err := redisClient.Incr(ctx, key).Result()
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, "redis error")
			c.Abort()
			return
		}

		if count == 1 {
			redisClient.Expire(ctx, key, window)
		}

		if count > int64(limit) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "global rate limit exceeded",
			})
			return
		}

		c.Next()
	}
}

func RateLimitIPMiddleware(redisClient *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := fmt.Sprintf("rate_limit:ip:%s", ip)

		count, err := redisClient.Incr(ctx, key).Result()
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, "redis error")
			c.Abort()
			return
		}

		if count == 1 {
			redisClient.Expire(ctx, key, window)
		}

		if count > int64(limit) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded (IP)",
			})
			return
		}

		c.Next()
	}
}
