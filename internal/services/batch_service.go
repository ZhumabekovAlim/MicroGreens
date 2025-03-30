package services

import (
	"MicroGreens/internal/models"
	"MicroGreens/internal/repositories"
	"context"
)

type BatchService struct {
	Repo *repositories.BatchRepository
}

func (s *BatchService) Create(ctx context.Context, b models.Batch) (models.Batch, error) {
	return s.Repo.Create(ctx, b)
}

func (s *BatchService) GetAll(ctx context.Context) ([]models.Batch, error) {
	return s.Repo.GetAll(ctx)
}

func (s *BatchService) GetByID(ctx context.Context, id int) (models.Batch, error) {
	return s.Repo.GetByID(ctx, id)
}

func (s *BatchService) Update(ctx context.Context, b models.Batch) error {
	return s.Repo.Update(ctx, b)
}

func (s *BatchService) Delete(ctx context.Context, id int) error {
	return s.Repo.Delete(ctx, id)
}
