package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"koalbot_api/internal/pagination"
	"koalbot_api/internal/repository"
	"koalbot_api/internal/service"
)

type MasterPenggunaHandler struct {
	service *service.MasterPenggunaService
}

func NewMasterPenggunaHandler(service *service.MasterPenggunaService) *MasterPenggunaHandler {
	return &MasterPenggunaHandler{service: service}
}

type createMasterPenggunaRequest struct {
	IDPengguna int64   `json:"id_pengguna"`
	Telegram   *string `json:"telegram"`
	Jenis      *string `json:"jenis"`
	Active     *bool   `json:"active"`
}

type updateMasterPenggunaRequest struct {
	IDPengguna *int64  `json:"id_pengguna"`
	Telegram   *string `json:"telegram"`
	Jenis      *string `json:"jenis"`
	Active     *bool   `json:"active"`
}

type masterPenggunaItem struct {
	ID         int64      `json:"id"`
	UUID       string     `json:"uuid"`
	IDPengguna int64      `json:"id_pengguna"`
	Telegram   *string    `json:"telegram"`
	Jenis      string     `json:"jenis"`
	Active     bool       `json:"active"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}

func (h *MasterPenggunaHandler) Create(c *gin.Context) {
	var req createMasterPenggunaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "invalid_json"})
		return
	}
	if req.IDPengguna == 0 {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "id_pengguna_required"})
		return
	}

	jenis := "stockity"
	if req.Jenis != nil {
		jenis = *req.Jenis
	}
	if !isValidJenis(jenis) {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "invalid_jenis"})
		return
	}

	active := false
	if req.Active != nil {
		active = *req.Active
	}

	item, err := h.service.Create(c.Request.Context(), req.IDPengguna, req.Telegram, jenis, active)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "server_error"})
		return
	}

	c.JSON(http.StatusCreated, masterPenggunaItem{
		ID:         item.ID,
		UUID:       item.UUID,
		IDPengguna: item.IDPengguna,
		Telegram:   item.Telegram,
		Jenis:      item.Jenis,
		Active:     item.Active,
		CreatedAt:  item.CreatedAt,
		UpdatedAt:  item.UpdatedAt,
		DeletedAt:  item.DeletedAt,
	})
}

func (h *MasterPenggunaHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "invalid_id"})
		return
	}

	item, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse{Error: "not_found"})
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "server_error"})
		return
	}

	c.JSON(http.StatusOK, masterPenggunaItem{
		ID:         item.ID,
		UUID:       item.UUID,
		IDPengguna: item.IDPengguna,
		Telegram:   item.Telegram,
		Jenis:      item.Jenis,
		Active:     item.Active,
		CreatedAt:  item.CreatedAt,
		UpdatedAt:  item.UpdatedAt,
		DeletedAt:  item.DeletedAt,
	})
}

func (h *MasterPenggunaHandler) List(c *gin.Context) {
	params, err := pagination.Parse(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "invalid_pagination"})
		return
	}
	jenis := c.Query("jenis")
	if jenis != "" && !isValidJenis(jenis) {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "invalid_jenis"})
		return
	}

	items, total, err := h.service.List(c.Request.Context(), params.Search, jenis, params.Limit, params.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "server_error"})
		return
	}

	responseItems := make([]masterPenggunaItem, 0, len(items))
	for _, item := range items {
		responseItems = append(responseItems, masterPenggunaItem{
			ID:         item.ID,
			UUID:       item.UUID,
			IDPengguna: item.IDPengguna,
			Telegram:   item.Telegram,
			Jenis:      item.Jenis,
			Active:     item.Active,
			CreatedAt:  item.CreatedAt,
			UpdatedAt:  item.UpdatedAt,
			DeletedAt:  item.DeletedAt,
		})
	}

	c.JSON(http.StatusOK, pagination.NewResponse(responseItems, total, params))
}

func (h *MasterPenggunaHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "invalid_id"})
		return
	}

	var req updateMasterPenggunaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "invalid_json"})
		return
	}

	update := repository.UpdateMasterPenggunaRequest{
		IDPengguna: req.IDPengguna,
		Telegram:   req.Telegram,
		Jenis:      req.Jenis,
		Active:     req.Active,
	}
	if req.Jenis != nil && !isValidJenis(*req.Jenis) {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "invalid_jenis"})
		return
	}

	if err := h.service.Update(c.Request.Context(), id, update); err != nil {
		switch {
		case errors.Is(err, service.ErrMasterPenggunaNoFields):
			c.JSON(http.StatusBadRequest, errorResponse{Error: "no_fields_to_update"})
		default:
			c.JSON(http.StatusInternalServerError, errorResponse{Error: "server_error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *MasterPenggunaHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "invalid_id"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "server_error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func isValidJenis(val string) bool {
	return val == "stockity" || val == "binomo" || val == "olymptrade"
}

func (h *MasterPenggunaHandler) Summary(c *gin.Context) {
	activeCount, inactiveCount, err := h.service.Summary(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "server_error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":    activeCount + inactiveCount,
		"active":   activeCount,
		"inactive": inactiveCount,
	})
}
