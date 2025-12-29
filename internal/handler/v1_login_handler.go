package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"koalbot_api/internal/model"
	"koalbot_api/internal/repository"
	"koalbot_api/internal/stockity"
)

const beTokenTTL = 72 * time.Hour

type V1LoginHandler struct {
	stockity *stockity.Client
	master   *repository.MasterPenggunaRepository
	detail   *repository.PenggunaDetailRepository
	secret   []byte
	apiURL   string
}

func NewV1LoginHandler(stockityClient *stockity.Client, master *repository.MasterPenggunaRepository, detail *repository.PenggunaDetailRepository, secret string, apiURL string) *V1LoginHandler {
	return &V1LoginHandler{
		stockity: stockityClient,
		master:   master,
		detail:   detail,
		secret:   []byte(secret),
		apiURL:   apiURL,
	}
}

type v1LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type v1LoginResponse struct {
	Token       string          `json:"token"`
	TokenAPI    string          `json:"token_api"`
	UserProfile penggunaProfile `json:"user_profile"`
	APIURL      string          `json:"api_url"`
}

type penggunaProfile struct {
	ID                         int64           `json:"id"`
	Avatar                     *string         `json:"avatar"`
	FirstName                  string          `json:"first_name"`
	LastName                   string          `json:"last_name"`
	Nickname                   string          `json:"nickname"`
	Balance                    float64         `json:"balance"`
	BalanceVersion             int64           `json:"balance_version"`
	Bonus                      float64         `json:"bonus"`
	Gender                     string          `json:"gender"`
	Email                      string          `json:"email"`
	EmailVerified              bool            `json:"email_verified"`
	Phone                      string          `json:"phone"`
	PhoneVerified              bool            `json:"phone_verified"`
	PhonePrefix                string          `json:"phone_prefix"`
	ReceiveNews                bool            `json:"receive_news"`
	ReceiveSMS                 bool            `json:"receive_sms"`
	ReceiveNotification        bool            `json:"receive_notification"`
	Country                    string          `json:"country"`
	CountryName                string          `json:"country_name"`
	Currency                   string          `json:"currency"`
	Birthday                   string          `json:"birthday"`
	Activate                   bool            `json:"activate"`
	PasswordIsSet              bool            `json:"password_is_set"`
	Tutorial                   bool            `json:"tutorial"`
	Coupons                    json.RawMessage `json:"coupons"`
	FreeDeals                  json.RawMessage `json:"free_deals"`
	Blocked                    bool            `json:"blocked"`
	AgreeRisk                  bool            `json:"agree_risk"`
	Agreed                     bool            `json:"agreed"`
	StatusGroup                string          `json:"status_group"`
	DocsVerified               bool            `json:"docs_verified"`
	RegisteredAt               *time.Time      `json:"registered_at"`
	StatusByDeposit            string          `json:"status_by_deposit"`
	StatusID                   int             `json:"status_id"`
	DepositsSum                float64         `json:"deposits_sum"`
	PushNotificationCategories json.RawMessage `json:"push_notification_categories"`
	PreserveName               bool            `json:"preserve_name"`
	RegistrationCountryISO     string          `json:"registration_country_iso"`
}

func (h *V1LoginHandler) Login(c *gin.Context) {
	deviceID := strings.TrimSpace(c.GetHeader("Device-Id"))
	deviceType := strings.TrimSpace(c.GetHeader("Device-Type"))
	if deviceID == "" || deviceType == "" {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "device_headers_required"})
		return
	}

	var req v1LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "invalid_json"})
		return
	}
	if req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "email_password_required"})
		return
	}

	signIn, err := h.stockity.SignIn(c.Request.Context(), deviceID, deviceType, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, stockity.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, errorResponse{Error: "invalid_credentials"})
			return
		}
		var upstreamErr *stockity.UpstreamError
		if errors.As(err, &upstreamErr) {
			if upstreamErr.Status == http.StatusUnprocessableEntity {
				c.JSON(http.StatusBadRequest, gin.H{"error": "upstream_validation", "message": upstreamErr.Body})
				return
			}
			c.JSON(http.StatusBadGateway, gin.H{"error": "upstream_error", "message": upstreamErr.Body})
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{"error": "upstream_error", "message": err.Error()})
		return
	}

	userID, err := stockity.ParseUserID(signIn.UserID)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "upstream_error", "message": "invalid user_id from upstream"})
		return
	}

	master, err := h.master.GetByIDPengguna(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, createErr := h.master.Create(c.Request.Context(), userID, nil, "stockity", false)
			if createErr != nil {
				c.JSON(http.StatusInternalServerError, errorResponse{Error: "server_error"})
				return
			}
			c.JSON(http.StatusForbidden, errorResponse{Error: "account_inactive"})
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "server_error"})
		return
	}

	if !master.Active {
		c.JSON(http.StatusForbidden, errorResponse{Error: "account_inactive"})
		return
	}

	detail, err := h.detail.GetByPenggunaID(c.Request.Context(), master.ID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusInternalServerError, errorResponse{Error: "server_error"})
			return
		}

		profile, err := h.stockity.GetProfile(c.Request.Context(), deviceID, deviceType, signIn.AuthToken)
		if err != nil {
			var upstreamErr *stockity.UpstreamError
			if errors.As(err, &upstreamErr) {
				c.JSON(http.StatusBadGateway, gin.H{"error": "profile_fetch_failed", "message": upstreamErr.Body})
				return
			}
			c.JSON(http.StatusBadGateway, errorResponse{Error: "profile_fetch_failed"})
			return
		}

		mapped, err := mapProfileToDetail(profile, master.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse{Error: "server_error"})
			return
		}

		if err := h.detail.Upsert(c.Request.Context(), mapped); err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse{Error: "server_error"})
			return
		}
		detail = mapped
	}

	token, err := h.issueBEToken(master, time.Now().UTC())
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "server_error"})
		return
	}

	c.JSON(http.StatusOK, v1LoginResponse{
		Token:       token,
		TokenAPI:    signIn.AuthToken,
		UserProfile: detailToProfile(detail),
		APIURL:      h.apiURL,
	})
}

func (h *V1LoginHandler) issueBEToken(master model.MasterPengguna, lastSeen time.Time) (string, error) {
	claims := jwt.MapClaims{
		"sub":       master.UUID,
		"user_id":   master.IDPengguna,
		"last_seen": lastSeen.Unix(),
		"exp":       time.Now().Add(beTokenTTL).Unix(),
		"iat":       time.Now().Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(h.secret)
}

func mapProfileToDetail(profile stockity.Profile, penggunaID int64) (model.PenggunaDetail, error) {
	var registeredAt *time.Time
	if profile.RegisteredAt != "" {
		parsed, err := time.Parse(time.RFC3339, profile.RegisteredAt)
		if err == nil {
			registeredAt = &parsed
		}
	}

	return model.PenggunaDetail{
		ID:                         profile.ID,
		PenggunaID:                 penggunaID,
		Avatar:                     profile.Avatar,
		FirstName:                  profile.FirstName,
		LastName:                   profile.LastName,
		Nickname:                   profile.Nickname,
		Balance:                    profile.Balance,
		BalanceVersion:             profile.BalanceVersion,
		Bonus:                      profile.Bonus,
		Gender:                     profile.Gender,
		Email:                      profile.Email,
		EmailVerified:              profile.EmailVerified,
		Phone:                      profile.Phone,
		PhoneVerified:              profile.PhoneVerified,
		PhonePrefix:                profile.PhonePrefix,
		ReceiveNews:                profile.ReceiveNews,
		ReceiveSMS:                 profile.ReceiveSMS,
		ReceiveNotification:        profile.ReceiveNotification,
		Country:                    profile.Country,
		CountryName:                profile.CountryName,
		Currency:                   profile.Currency,
		Birthday:                   profile.Birthday,
		Activate:                   profile.Activate,
		PasswordIsSet:              profile.PasswordIsSet,
		Tutorial:                   profile.Tutorial,
		Coupons:                    profile.Coupons,
		FreeDeals:                  profile.FreeDeals,
		Blocked:                    profile.Blocked,
		AgreeRisk:                  profile.AgreeRisk,
		Agreed:                     profile.Agreed,
		StatusGroup:                profile.StatusGroup,
		DocsVerified:               profile.DocsVerified,
		RegisteredAt:               registeredAt,
		StatusByDeposit:            profile.StatusByDeposit,
		StatusID:                   profile.StatusID,
		DepositsSum:                profile.DepositsSum,
		PushNotificationCategories: profile.PushNotificationCategories,
		PreserveName:               profile.PreserveName,
		RegistrationCountryISO:     profile.RegistrationCountryISO,
	}, nil
}

func detailToProfile(detail model.PenggunaDetail) penggunaProfile {
	return penggunaProfile{
		ID:                         detail.ID,
		Avatar:                     detail.Avatar,
		FirstName:                  detail.FirstName,
		LastName:                   detail.LastName,
		Nickname:                   detail.Nickname,
		Balance:                    detail.Balance,
		BalanceVersion:             detail.BalanceVersion,
		Bonus:                      detail.Bonus,
		Gender:                     detail.Gender,
		Email:                      detail.Email,
		EmailVerified:              detail.EmailVerified,
		Phone:                      detail.Phone,
		PhoneVerified:              detail.PhoneVerified,
		PhonePrefix:                detail.PhonePrefix,
		ReceiveNews:                detail.ReceiveNews,
		ReceiveSMS:                 detail.ReceiveSMS,
		ReceiveNotification:        detail.ReceiveNotification,
		Country:                    detail.Country,
		CountryName:                detail.CountryName,
		Currency:                   detail.Currency,
		Birthday:                   detail.Birthday,
		Activate:                   detail.Activate,
		PasswordIsSet:              detail.PasswordIsSet,
		Tutorial:                   detail.Tutorial,
		Coupons:                    detail.Coupons,
		FreeDeals:                  detail.FreeDeals,
		Blocked:                    detail.Blocked,
		AgreeRisk:                  detail.AgreeRisk,
		Agreed:                     detail.Agreed,
		StatusGroup:                detail.StatusGroup,
		DocsVerified:               detail.DocsVerified,
		RegisteredAt:               detail.RegisteredAt,
		StatusByDeposit:            detail.StatusByDeposit,
		StatusID:                   detail.StatusID,
		DepositsSum:                detail.DepositsSum,
		PushNotificationCategories: detail.PushNotificationCategories,
		PreserveName:               detail.PreserveName,
		RegistrationCountryISO:     detail.RegistrationCountryISO,
	}
}
