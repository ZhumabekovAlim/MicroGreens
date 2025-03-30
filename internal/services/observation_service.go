package services

import (
	"MicroGreens/internal/models"
	"MicroGreens/internal/repositories"
	"context"
)

type ObservationService struct {
	Repo *repositories.ObservationRepository
}

func (s *ObservationService) Create(ctx context.Context, o models.Observation) (models.Observation, error) {
	return s.Repo.Create(ctx, o)
}

func (s *ObservationService) GetAll(ctx context.Context) ([]models.Observation, error) {
	return s.Repo.GetAll(ctx)
}

func (s *ObservationService) GetByID(ctx context.Context, id int) (models.Observation, error) {
	return s.Repo.GetByID(ctx, id)
}

func (s *ObservationService) Update(ctx context.Context, o models.Observation) error {
	return s.Repo.Update(ctx, o)
}

func (s *ObservationService) Delete(ctx context.Context, id int) error {
	return s.Repo.Delete(ctx, id)
}
