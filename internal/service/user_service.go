package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"koalbot_api/internal/model"
	"koalbot_api/internal/repository"
)

var ErrInvalidRole = errors.New("invalid_role")
var ErrNoFieldsToUpdate = errors.New("no_fields_to_update")

type UserService struct {
	users *repository.UserRepository
}

func NewUserService(users *repository.UserRepository) *UserService {
	return &UserService{users: users}
}

func (s *UserService) Register(ctx context.Context, username, password, role, createdBy string) (string, string, error) {
	if !isValidRole(role) {
		return "", "", ErrInvalidRole
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	uid, storedRole, err := s.users.CreateUser(ctx, username, string(hash), role, createdBy)
	if err != nil {
		return "", "", err
	}

	return uid, storedRole, nil
}

func (s *UserService) Update(ctx context.Context, uid string, req repository.UpdateUserRequest) error {
	if req.Role != nil && !isValidRole(*req.Role) {
		return ErrInvalidRole
	}

	if req.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		req.Password = ptrString(string(hash))
	}

	if req.Username == nil && req.Password == nil && req.Role == nil && req.Active == nil {
		return ErrNoFieldsToUpdate
	}

	return s.users.UpdateUser(ctx, uid, req)
}

func (s *UserService) Delete(ctx context.Context, uid, deletedBy string) error {
	return s.users.SoftDeleteUser(ctx, uid, deletedBy)
}

func (s *UserService) List(ctx context.Context, search string, limit, offset int) ([]model.User, int, error) {
	return s.users.ListUsers(ctx, search, limit, offset)
}

func isValidRole(role string) bool {
	return role == "admin" || role == "viewer"
}

func ptrString(val string) *string {
	return &val
}
