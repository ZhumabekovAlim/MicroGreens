package services

import (
	"MicroGreens/internal/models"
	"MicroGreens/internal/repositories"
	"context"
)

type MicrogreenService struct {
	Repo *repositories.MicrogreenRepository
}

func (s *MicrogreenService) Create(ctx context.Context, m models.Microgreen) (models.Microgreen, error) {
	return s.Repo.Create(ctx, m)
}

func (s *MicrogreenService) GetAll(ctx context.Context) ([]models.Microgreen, error) {
	return s.Repo.GetAll(ctx)
}

func (s *MicrogreenService) GetByID(ctx context.Context, id int) (models.Microgreen, error) {
	return s.Repo.GetByID(ctx, id)
}

func (s *MicrogreenService) Update(ctx context.Context, m models.Microgreen) error {
	return s.Repo.Update(ctx, m)
}

func (s *MicrogreenService) Delete(ctx context.Context, id int) error {
	return s.Repo.Delete(ctx, id)
}
