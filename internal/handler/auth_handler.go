package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"koalbot_api/internal/service"
)

type AuthHandler struct {
	auth   *service.AuthService
	tokens *service.TokenService
}

func NewAuthHandler(auth *service.AuthService, tokens *service.TokenService) *AuthHandler {
	return &AuthHandler{auth: auth, tokens: tokens}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token            string    `json:"token"`
	RefreshToken     string    `json:"refresh_token"`
	ExpiresAt        time.Time `json:"expires_at"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
	User             userInfo  `json:"user"`
}

type userInfo struct {
	UID      string     `json:"uid"`
	Username string     `json:"username"`
	Role     string     `json:"role"`
	LastSeen *time.Time `json:"last_seen"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "invalid_json"})
		return
	}
	if req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "username_and_password_required"})
		return
	}

	user, lastSeen, err := h.auth.Authenticate(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, errorResponse{Error: "invalid_credentials"})
		case errors.Is(err, service.ErrUserInactive):
			c.JSON(http.StatusUnauthorized, errorResponse{Error: "user_inactive"})
		default:
			c.JSON(http.StatusInternalServerError, errorResponse{Error: "server_error"})
		}
		return
	}

	accessToken, refreshToken, accessExp, refreshExp, err := h.tokens.IssueTokens(c.Request.Context(), user, lastSeen)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "server_error"})
		return
	}

	c.JSON(http.StatusOK, loginResponse{
		Token:            accessToken,
		RefreshToken:     refreshToken,
		ExpiresAt:        accessExp,
		RefreshExpiresAt: refreshExp,
		User: userInfo{
			UID:      user.UID,
			Username: user.Username,
			Role:     user.Role,
			LastSeen: user.LastSeen,
		},
	})
}
