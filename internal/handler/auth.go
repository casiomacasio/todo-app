package handler

import (
	"errors"
	"net/http"
	"github.com/casiomacasio/todo-app/internal/domain"
	"github.com/casiomacasio/todo-app/internal/repository"
	"github.com/gin-gonic/gin"
)

func (h *Handler) signUp(c *gin.Context) {
	var input domain.User
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid body") 
		return
	}

	id, err := h.service.Authorization.CreateUser(input)
	if err != nil {
		if errors.Is(err, repository.ErrUsernameExists) {
			newErrorResponse(c, http.StatusConflict, "Username already in use")
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.GetUser(input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	refreshToken, err := h.service.GenerateRefreshToken(user.Id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	accessToken, err := h.service.GenerateToken(user.Id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.SetCookie("refresh_token", refreshToken, 30*24*60*60, "/", "", true, true) 
	c.SetCookie("access_token", accessToken, 15*60, "/", "", true, true)      

	c.JSON(http.StatusOK, map[string]string{
		"message": "logged in successfully",
	})
}


func (h *Handler) refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "no refresh token cookie")
		return
	}

	userId, err := h.service.GetUserByRefreshToken(refreshToken)
	if err != nil {
		if errors.Is(err, repository.ErrRefreshTokenExpired) {
			c.Header("RefreshToken-Expired", "true")
			newErrorResponse(c, http.StatusUnauthorized, "refresh_token expired, must re-login")
			return
		}
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	newAccessToken, err := h.service.GenerateToken(userId)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	newRefreshToken, err := h.service.Authorization.GenerateRefreshToken(userId)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "failed to generate new refresh token")
		return
	}

	c.SetCookie("refresh_token", newRefreshToken, 30*24*60*60, "/", "", true, true)
	c.SetCookie("access_token", newAccessToken, 60*60, "/", "", true, true)

	c.JSON(http.StatusOK, map[string]string{
		"message": "token refreshed",
	})
}


func (h *Handler) logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err == nil {
		_ = h.service.Authorization.RevokeRefreshToken(refreshToken) 
	}
	c.SetCookie("refresh_token", "", -1, "/", "", true, true)
	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}