package services

import (
	"MicroGreens/internal/models"
	"MicroGreens/internal/repositories"
	"context"
)

type AdviceService struct {
	Repo *repositories.AdviceRepository
}

func (s *AdviceService) Create(ctx context.Context, a models.AdviceMessage) (models.AdviceMessage, error) {
	return s.Repo.Create(ctx, a)
}

func (s *AdviceService) GetAll(ctx context.Context) ([]models.AdviceMessage, error) {
	return s.Repo.GetAll(ctx)
}

func (s *AdviceService) GetByID(ctx context.Context, id int) (models.AdviceMessage, error) {
	return s.Repo.GetByID(ctx, id)
}

func (s *AdviceService) Update(ctx context.Context, a models.AdviceMessage) error {
	return s.Repo.Update(ctx, a)
}

func (s *AdviceService) Delete(ctx context.Context, id int) error {
	return s.Repo.Delete(ctx, id)
}
