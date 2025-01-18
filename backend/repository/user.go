package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/yDog-1/wodun/backend/generated/dbstore"
	"github.com/yDog-1/wodun/backend/graph/model"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) GetUser(ctx context.Context, uniqueName string) (*model.User, error) {
	query := dbstore.New(r.db)
	user, err := query.GetUser(ctx, uniqueName)
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:          fmt.Sprint(user.ID),
		UniqueName:  user.UniqueName,
		DisplayName: user.DisplayName,
		Email:       user.Email,
	}, nil
}

func (r *userRepository) CreateUser(ctx context.Context, input *model.CreateUserInput) (string, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()
	query := dbstore.New(tx)

	err = query.CreateUser(ctx, dbstore.CreateUserParams{
		UniqueName:  input.UniqueName,
		DisplayName: input.DisplayName,
		Email:       input.Email,
	})
	if err != nil {
		return "", err
	}

	id, err := query.LastInsertId(ctx)
	if err != nil {
		return "", err
	}
	err = tx.Commit()
	if err != nil {
		return "", err
	}
	return fmt.Sprint(id), nil
}

func (r *userRepository) UpdateUser(ctx context.Context, id string, input *model.UpdateUserInput) error {
	query := dbstore.New(r.db)

	var un, dn, em sql.NullString
	if input.UniqueName != nil {
		un.String = *input.UniqueName
		un.Valid = true
	}
	if input.DisplayName != nil {
		dn.String = *input.DisplayName
		dn.Valid = true
	}
	if input.Email != nil {
		em.String = *input.Email
		em.Valid = true
	}
	uintID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return err
	}

	err = query.UpdateUser(ctx, dbstore.UpdateUserParams{
		UniqueName:  un,
		DisplayName: dn,
		Email:       em,
		ID:          uintID,
	})

	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) DeleteUser(ctx context.Context, uniqueName string) error {
	query := dbstore.New(r.db)
	err := query.DeleteUser(ctx, uniqueName)
	if err != nil {
		return err
	}
	return nil
}
