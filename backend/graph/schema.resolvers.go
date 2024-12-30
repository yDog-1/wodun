package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.61

import (
	"context"
	"fmt"

	"github.com/yDog-1/wodun/backend/graph/model"
)

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUserInput) (*model.AuthPayload, error) {
	panic(fmt.Errorf("not implemented: CreateUser - createUser"))
}

// UpdateUser is the resolver for the updateUser field.
func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input model.UpdateUserInput) (bool, error) {
	panic(fmt.Errorf("not implemented: UpdateUser - updateUser"))
}

// SendMagicLink is the resolver for the sendMagicLink field.
func (r *mutationResolver) SendMagicLink(ctx context.Context, email string) (bool, error) {
	panic(fmt.Errorf("not implemented: SendMagicLink - sendMagicLink"))
}

// VerifyMagicLink is the resolver for the verifyMagicLink field.
func (r *mutationResolver) VerifyMagicLink(ctx context.Context, token string) (*model.AuthPayload, error) {
	panic(fmt.Errorf("not implemented: VerifyMagicLink - verifyMagicLink"))
}

// RefreshToken is the resolver for the refreshToken field.
func (r *mutationResolver) RefreshToken(ctx context.Context, refreshToken string) (*model.AuthPayload, error) {
	panic(fmt.Errorf("not implemented: RefreshToken - refreshToken"))
}

// Me is the resolver for the me field.
func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	panic(fmt.Errorf("not implemented: Me - me"))
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	panic(fmt.Errorf("not implemented: User - user"))
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }