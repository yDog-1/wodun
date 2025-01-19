package service

import (
	"context"

	"github.com/yDog-1/wodun/backend/graph/model"
)

type userRepository interface {
	GetUser(ctx context.Context, uniqueName string) (*model.User, error)
	CreateUser(ctx context.Context, input *model.CreateUserInput) (string, error)
	UpdateUser(ctx context.Context, id string, input *model.UpdateUserInput) error
	DeleteUser(ctx context.Context, uniqueName string) error
	// ListUsers(ctx context.Context) ([]*model.User, error)
}
type UserService struct {
	repo userRepository
}

func NewUserService(repo userRepository) *UserService {
	return &UserService{repo}
}

func (s *UserService) GetUser(ctx context.Context, uniqueName string) (*model.User, error) {
	return s.repo.GetUser(ctx, uniqueName)
}

func (s *UserService) CreateUser(ctx context.Context, input *model.CreateUserInput) (string, error) {
	return s.repo.CreateUser(ctx, input)
}

func (s *UserService) UpdateUser(ctx context.Context, id string, input *model.UpdateUserInput) error {
	return s.repo.UpdateUser(ctx, id,
		&model.UpdateUserInput{
			ID:          id,
			UniqueName:  input.UniqueName,
			DisplayName: input.DisplayName,
			Email:       input.Email,
		},
	)
}
func (s *UserService) DeleteUser(ctx context.Context, uniqueName string) error {
	return s.repo.DeleteUser(ctx, uniqueName)
}
