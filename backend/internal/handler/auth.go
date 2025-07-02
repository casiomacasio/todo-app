package handler

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/casiomacasio/todo-app/backend/internal/domain"
	"github.com/casiomacasio/todo-app/backend/internal/repository"
	"github.com/gin-gonic/gin"
)


// @Summary User registration
// @Tags auth
// @Accept json
// @Produce json
// @Param input body domain.CreateUserRequest true "User credentials"
// @Success 200 {object} map[string]interface{} "Registered user ID"
// @Failure 400 {object} errorResponse "Invalid request body"
// @Failure 409 {object} errorResponse "Username already exists"
// @Failure 500 {object} errorResponse "Server error"
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var input domain.CreateUserRequest
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

// @Summary User login
// @Tags auth
// @Accept json
// @Produce json
// @Param input body signInInput true "Username and password"
// @Success 200 {object} map[string]string "Login success message; sets cookies: access_token, refresh_token, refresh_token_id"
// @Failure 400 {object} errorResponse "Invalid request or credentials"
// @Failure 500 {object} errorResponse "Server error"
// @Router /auth/sign-in [post]
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

	tokenID, refreshToken, err := h.service.GenerateRefreshToken(user.Id)
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
	c.SetCookie("refresh_token_id", tokenID, 30*24*60*60, "/", "", true, true)
	c.SetCookie("access_token", accessToken, 15*60, "/", "", true, true)

	c.JSON(http.StatusOK, map[string]string{
		"message": "logged in successfully",
	})
}

// @Summary Refresh access token
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "New access and refresh tokens are set in cookies"
// @Failure 400 {object} errorResponse "Invalid token format"
// @Failure 401 {object} errorResponse "Missing or expired refresh token"
// @Failure 500 {object} errorResponse "Server error"
// @Router /auth/refresh [post]
func (h *Handler) refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "no refresh token cookie")
		return
	}
	refreshTokenId, err := c.Cookie("refresh_token_id")
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "refresh_token_id is missed")
		return
	}
	refreshTokenUUID, err := uuid.Parse(refreshTokenId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid refresh token format")
		return
	}
	userId, err := h.service.GetUserByRefreshTokenAndRefreshTokenId(refreshToken, refreshTokenUUID)
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

	newTokenId, newRefreshToken, err := h.service.Authorization.GenerateRefreshToken(userId)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "failed to generate new refresh token")
		return
	}

	c.SetCookie("refresh_token", newRefreshToken, 30*24*60*60, "/", "", true, true)
	c.SetCookie("refresh_token_id", newTokenId, 30*24*60*60, "/", "", true, true)
	c.SetCookie("access_token", newAccessToken, 60*60, "/", "", true, true)

	c.JSON(http.StatusOK, map[string]string{
		"message": "token refreshed",
	})
}

// @Summary User logout
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Logout success message; cookies cleared"
// @Failure 400 {object} errorResponse "Invalid token format"
// @Failure 401 {object} errorResponse "Missing or expired refresh token"
// @Failure 500 {object} errorResponse "Server error"
// @Router /auth/logout [post]
func (h *Handler) logout(c *gin.Context) {
	refreshTokenId, err := c.Cookie("refresh_token_id")
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "refresh_token_id is missed")
		return
	}
	refreshTokenUUID, err := uuid.Parse(refreshTokenId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid refresh token format")
		return
	}
	if err := h.service.Authorization.RevokeRefreshToken(refreshTokenUUID); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "couldn't revoke refresh token")
	}

	c.SetCookie("refresh_token", "", -1, "/", "", true, true)
	c.SetCookie("refresh_token_id", "", -1, "/", "", true, true)
	c.SetCookie("access_token", "", -1, "/", "", true, true) 

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}