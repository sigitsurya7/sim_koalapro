package service

import (
	"context"
	"errors"

	"koalbot_api/internal/model"
	"koalbot_api/internal/repository"
)

var ErrMasterPenggunaNoFields = errors.New("no_fields_to_update")

type MasterPenggunaService struct {
	repo *repository.MasterPenggunaRepository
}

func NewMasterPenggunaService(repo *repository.MasterPenggunaRepository) *MasterPenggunaService {
	return &MasterPenggunaService{repo: repo}
}

func (s *MasterPenggunaService) Create(ctx context.Context, idPengguna int64, telegram *string, jenis string, active bool) (model.MasterPengguna, error) {
	return s.repo.Create(ctx, idPengguna, telegram, jenis, active)
}

func (s *MasterPenggunaService) GetByID(ctx context.Context, id int64) (model.MasterPengguna, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *MasterPenggunaService) List(ctx context.Context, search string, jenis string, limit, offset int) ([]model.MasterPengguna, int, error) {
	return s.repo.List(ctx, search, jenis, limit, offset)
}

func (s *MasterPenggunaService) Update(ctx context.Context, id int64, req repository.UpdateMasterPenggunaRequest) error {
	if req.IDPengguna == nil && req.Telegram == nil && req.Jenis == nil && req.Active == nil {
		return ErrMasterPenggunaNoFields
	}
	return s.repo.Update(ctx, id, req.IDPengguna, req.Telegram, req.Jenis, req.Active)
}

func (s *MasterPenggunaService) Delete(ctx context.Context, id int64) error {
	return s.repo.SoftDelete(ctx, id)
}

func (s *MasterPenggunaService) Summary(ctx context.Context) (int, int, error) {
	return s.repo.CountByActive(ctx)
}
