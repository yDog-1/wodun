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
type userService struct {
	repo userRepository
}

func NewUserService(repo userRepository) *userService {
	return &userService{repo}
}

func (s *userService) GetUser(ctx context.Context, uniqueName string) (*model.User, error) {
	return s.repo.GetUser(ctx, uniqueName)
}

func (s *userService) CreateUser(ctx context.Context, input *model.CreateUserInput) (string, error) {
	return s.repo.CreateUser(ctx, input)
}

func (s *userService) UpdateUser(ctx context.Context, id string, input *model.UpdateUserInput) error {
	return s.repo.UpdateUser(ctx, id,
		&model.UpdateUserInput{
			ID:          id,
			UniqueName:  input.UniqueName,
			DisplayName: input.DisplayName,
			Email:       input.Email,
		},
	)
}
func (s *userService) DeleteUser(ctx context.Context, uniqueName string) error {
	return s.repo.DeleteUser(ctx, uniqueName)
}
