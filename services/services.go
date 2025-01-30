package services

import (
    "context"
    "fx.dependency.injection/repositories"
    "fx.dependency.injection/models"
)

type UserService struct {
    repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
    return s.repo.CreateUser(ctx, user)
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
    return s.repo.GetUserByID(ctx, id)
}

func (s *UserService) UpdateUser(ctx context.Context, user *models.User) error {
    return s.repo.UpdateUser(ctx, user)
}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
    return s.repo.DeleteUser(ctx, id)
}