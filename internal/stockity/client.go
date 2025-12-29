package stockity

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var ErrInvalidCredentials = errors.New("invalid_credentials")

type UpstreamError struct {
	Status   int
	Endpoint string
	Body     string
}

func (e *UpstreamError) Error() string {
	return fmt.Sprintf("stockity %s status %d: %s", e.Endpoint, e.Status, e.Body)
}

type Client struct {
	baseURL string
	http    *http.Client
}

func NewClient(baseURL string, timeout time.Duration) *Client {
	trimmed := strings.TrimRight(baseURL, "/")
	return &Client{
		baseURL: trimmed,
		http: &http.Client{
			Timeout: timeout,
		},
	}
}

type SignInResponse struct {
	AuthToken string `json:"authtoken"`
	UserID    string `json:"user_id"`
}

type signInWrapper struct {
	Data SignInResponse `json:"data"`
}

type Profile struct {
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
	RegisteredAt               string          `json:"registered_at"`
	StatusByDeposit            string          `json:"status_by_deposit"`
	StatusID                   int             `json:"status_id"`
	DepositsSum                float64         `json:"deposits_sum"`
	PushNotificationCategories json.RawMessage `json:"push_notification_categories"`
	PreserveName               bool            `json:"preserve_name"`
	RegistrationCountryISO     string          `json:"registration_country_iso"`
}

type profileResponse struct {
	Data *Profile `json:"data"`
}

func (c *Client) SignIn(ctx context.Context, deviceID, deviceType, email, password string) (SignInResponse, error) {
	payload := map[string]string{
		"email":    email,
		"password": password,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return SignInResponse{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/passport/v2/sign_in", bytes.NewReader(body))
	if err != nil {
		return SignInResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Device-Id", deviceID)
	req.Header.Set("Device-Type", deviceType)

	resp, err := c.http.Do(req)
	if err != nil {
		return SignInResponse{}, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return SignInResponse{}, err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return SignInResponse{}, ErrInvalidCredentials
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return SignInResponse{}, &UpstreamError{
			Status:   resp.StatusCode,
			Endpoint: "sign_in",
			Body:     limitBody(respBody),
		}
	}

	var wrapped signInWrapper
	if err := json.Unmarshal(respBody, &wrapped); err == nil && wrapped.Data.AuthToken != "" {
		if wrapped.Data.UserID == "" {
			return SignInResponse{}, errors.New("invalid sign_in response")
		}
		return wrapped.Data, nil
	}

	var result SignInResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return SignInResponse{}, err
	}
	if result.AuthToken == "" || result.UserID == "" {
		return SignInResponse{}, errors.New("invalid sign_in response")
	}

	return result, nil
}

func (c *Client) GetProfile(ctx context.Context, deviceID, deviceType, authToken string) (Profile, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/platform/private/v2/profile", nil)
	if err != nil {
		return Profile{}, err
	}
	req.Header.Set("Device-Id", deviceID)
	req.Header.Set("Device-Type", deviceType)
	req.Header.Set("Authorization-Token", authToken)

	resp, err := c.http.Do(req)
	if err != nil {
		return Profile{}, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Profile{}, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return Profile{}, &UpstreamError{
			Status:   resp.StatusCode,
			Endpoint: "profile",
			Body:     limitBody(respBody),
		}
	}

	var wrapper profileResponse
	if err := json.Unmarshal(respBody, &wrapper); err == nil && wrapper.Data != nil {
		return *wrapper.Data, nil
	}

	var profile Profile
	if err := json.Unmarshal(respBody, &profile); err != nil {
		return Profile{}, err
	}
	if profile.ID == 0 {
		return Profile{}, errors.New("invalid profile response")
	}

	return profile, nil
}

func limitBody(body []byte) string {
	const maxLen = 1024
	if len(body) == 0 {
		return ""
	}
	if len(body) > maxLen {
		return string(body[:maxLen]) + "...(truncated)"
	}
	return string(body)
}
