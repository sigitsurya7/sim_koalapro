package service

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"koalbot_api/internal/model"
	"koalbot_api/internal/repository"
)

const (
	accessTokenTTL  = 24 * time.Hour
	refreshTokenTTL = 7 * 24 * time.Hour
)

type TokenService struct {
	secret []byte
	repo   *repository.TokenRepository
}

type tokenClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	LastSeen int64  `json:"last_seen"`
	jwt.RegisteredClaims
}

func NewTokenService(secret string, repo *repository.TokenRepository) *TokenService {
	return &TokenService{
		secret: []byte(secret),
		repo:   repo,
	}
}

func (s *TokenService) IssueTokens(ctx context.Context, user model.User, lastSeen time.Time) (string, string, time.Time, time.Time, error) {
	accessExp := time.Now().Add(accessTokenTTL)
	refreshExp := time.Now().Add(refreshTokenTTL)

	accessToken, err := s.signToken(user, lastSeen, accessExp)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}

	refreshToken, err := s.signToken(user, lastSeen, refreshExp)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}

	if err := s.repo.SaveRefreshToken(ctx, user.UID, refreshToken, lastSeen, refreshExp); err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}

	return accessToken, refreshToken, accessExp, refreshExp, nil
}

func (s *TokenService) signToken(user model.User, lastSeen time.Time, expiresAt time.Time) (string, error) {
	claims := tokenClaims{
		Username: user.Username,
		Role:     user.Role,
		LastSeen: lastSeen.Unix(),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.UID,
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
}
