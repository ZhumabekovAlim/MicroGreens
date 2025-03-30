package services

import (
	"MicroGreens/internal/models"
	"MicroGreens/internal/repositories"
	"context"
)

type PhotoService struct {
	Repo *repositories.PhotoRepository
}

func (s *PhotoService) Create(ctx context.Context, p models.ObservationPhoto) (models.ObservationPhoto, error) {
	return s.Repo.Create(ctx, p)
}

func (s *PhotoService) GetAll(ctx context.Context) ([]models.ObservationPhoto, error) {
	return s.Repo.GetAll(ctx)
}

func (s *PhotoService) GetByID(ctx context.Context, id int) (models.ObservationPhoto, error) {
	return s.Repo.GetByID(ctx, id)
}

func (s *PhotoService) Update(ctx context.Context, p models.ObservationPhoto) error {
	return s.Repo.Update(ctx, p)
}

func (s *PhotoService) Delete(ctx context.Context, id int) error {
	return s.Repo.Delete(ctx, id)
}
