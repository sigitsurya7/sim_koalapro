package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"koalbot_api/internal/middleware"
	"koalbot_api/internal/pagination"
	"koalbot_api/internal/repository"
	"koalbot_api/internal/service"
)

type UserHandler struct {
	users *service.UserService
}

func NewUserHandler(users *service.UserService) *UserHandler {
	return &UserHandler{users: users}
}

type registerUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type updateUserRequest struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
	Role     *string `json:"role"`
	Active   *bool   `json:"active"`
}

type userItem struct {
	UID       string     `json:"uid"`
	Username  string     `json:"username"`
	Role      string     `json:"role"`
	Active    bool       `json:"active"`
	CreatedAt time.Time  `json:"created_at"`
	CreatedBy *string    `json:"created_by"`
	UpdatedAt *time.Time `json:"updated_at"`
	UpdatedBy *string    `json:"updated_by"`
}

type registerUserResponse struct {
	UID      string `json:"uid"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req registerUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "invalid_json"})
		return
	}
	if req.Username == "" || req.Password == "" || req.Role == "" {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "username_password_role_required"})
		return
	}

	authCtx, ok := middleware.GetAuthContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, errorResponse{Error: "unauthorized"})
		return
	}

	uid, role, err := h.users.Register(c.Request.Context(), req.Username, req.Password, req.Role, authCtx.UID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidRole):
			c.JSON(http.StatusBadRequest, errorResponse{Error: "invalid_role"})
		default:
			c.JSON(http.StatusInternalServerError, errorResponse{Error: "server_error"})
		}
		return
	}

	c.JSON(http.StatusCreated, registerUserResponse{UID: uid, Username: req.Username, Role: role})
}

func (h *UserHandler) List(c *gin.Context) {
	params, err := pagination.Parse(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "invalid_pagination"})
		return
	}

	users, total, err := h.users.List(c.Request.Context(), params.Search, params.Limit, params.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "server_error"})
		return
	}

	items := make([]userItem, 0, len(users))
	for _, user := range users {
		items = append(items, userItem{
			UID:       user.UID,
			Username:  user.Username,
			Role:      user.Role,
			Active:    user.Active,
			CreatedAt: user.CreatedAt,
			CreatedBy: user.CreatedBy,
			UpdatedAt: user.UpdatedAt,
			UpdatedBy: user.UpdatedBy,
		})
	}

	c.JSON(http.StatusOK, pagination.NewResponse(items, total, params))
}

func (h *UserHandler) Update(c *gin.Context) {
	uid := c.Param("uid")
	if uid == "" {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "uid_required"})
		return
	}

	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "invalid_json"})
		return
	}

	authCtx, ok := middleware.GetAuthContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, errorResponse{Error: "unauthorized"})
		return
	}

	update := repository.UpdateUserRequest{
		Username:  req.Username,
		Password:  req.Password,
		Role:      req.Role,
		Active:    req.Active,
		UpdatedBy: &authCtx.UID,
	}

	if err := h.users.Update(c.Request.Context(), uid, update); err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidRole):
			c.JSON(http.StatusBadRequest, errorResponse{Error: "invalid_role"})
		case errors.Is(err, service.ErrNoFieldsToUpdate):
			c.JSON(http.StatusBadRequest, errorResponse{Error: "no_fields_to_update"})
		default:
			c.JSON(http.StatusInternalServerError, errorResponse{Error: "server_error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *UserHandler) Delete(c *gin.Context) {
	uid := c.Param("uid")
	if uid == "" {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "uid_required"})
		return
	}
	authCtx, ok := middleware.GetAuthContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, errorResponse{Error: "unauthorized"})
		return
	}

	if err := h.users.Delete(c.Request.Context(), uid, authCtx.UID); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "server_error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
