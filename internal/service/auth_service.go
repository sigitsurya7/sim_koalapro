package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"koalbot_api/internal/model"
	"koalbot_api/internal/repository"
)

var ErrInvalidCredentials = errors.New("invalid_credentials")
var ErrUserInactive = errors.New("user_inactive")

type AuthService struct {
	users *repository.UserRepository
}

func NewAuthService(users *repository.UserRepository) *AuthService {
	return &AuthService{users: users}
}

func (s *AuthService) Authenticate(ctx context.Context, username, password string) (model.User, time.Time, error) {
	user, err := s.users.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, time.Time{}, ErrInvalidCredentials
		}
		return model.User{}, time.Time{}, err
	}

	if !user.Active || user.DeletedAt != nil {
		return model.User{}, time.Time{}, ErrUserInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return model.User{}, time.Time{}, ErrInvalidCredentials
	}

	lastSeen := time.Now().UTC()
	if err := s.users.UpdateLastSeen(ctx, user.UID, lastSeen); err != nil {
		return model.User{}, time.Time{}, err
	}
	user.LastSeen = &lastSeen

	return user, lastSeen, nil
}
