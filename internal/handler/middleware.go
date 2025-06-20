package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/casiomacasio/todo-app/internal/service"
)

const userCtx = "userId"

func (h Handler) userIdentity(c *gin.Context) {
	token, err := c.Cookie("access_token")
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "missing access token cookie")
		return
	}

	userID, err := h.service.Authorization.ParseToken(token)
	if err != nil {
		if errors.Is(err, service.ErrTokenExpired) {
			c.Header("Token-Expired", "true")
			newErrorResponse(c, http.StatusUnauthorized, "token expired, must be refreshed")
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
	idInt, ok := id.(int)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id is not of valid type")
		return 0, errors.New("user id is not of valid type")
	}
	return idInt, nil
}
