package services

import (
	"MicroGreens/internal/models"
	"MicroGreens/internal/repositories"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo *repositories.UserRepository
}

func (s *UserService) Create(ctx context.Context, u models.User) (models.User, error) {
	return s.Repo.Create(ctx, u)
}

func (s *UserService) GetAll(ctx context.Context) ([]models.User, error) {
	return s.Repo.GetAll(ctx)
}

func (s *UserService) GetByID(ctx context.Context, id int) (models.User, error) {
	return s.Repo.GetByID(ctx, id)
}

func (s *UserService) Update(ctx context.Context, u models.User) error {
	return s.Repo.Update(ctx, u)
}

func (s *UserService) Delete(ctx context.Context, id int) error {
	return s.Repo.Delete(ctx, id)
}

func (s *UserService) Authenticate(email, password string) (int, error) {
	user, err := s.Repo.GetUserByEmail(email)
	if err != nil {
		return 0, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return 0, errors.New("incorrect password")
	}

	return user.ID, nil
}
