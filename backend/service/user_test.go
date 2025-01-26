package service_test

import (
	"context"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yDog-1/wodun/backend/graph/model"
	"github.com/yDog-1/wodun/backend/pkg"
	"github.com/yDog-1/wodun/backend/pkg/testing/container"
	"github.com/yDog-1/wodun/backend/repository"
	"github.com/yDog-1/wodun/backend/service"
)

func Test_ユーザーを作成する(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	db, terminate := container.MysqlContainer(
		t,
		ctx,
		container.MySQLcontainerInput(),
	)
	defer terminate()

	repo := repository.NewUserRepository(db)
	s := service.NewUserService(repo)
	err := s.DeleteUser(ctx, "ydog")
	require.NoError(t, err)
	input := model.CreateUserInput{
		UniqueName:  "ydog",
		DisplayName: "yDog",
		Email:       "ydog@example.com",
	}

	id, err := s.CreateUser(ctx, &input)
	require.NoError(t, err)
	assert.NotEqual(t, "", id)

	user, err := s.GetUser(ctx, "ydog")
	assert.Nil(t, err)
	assert.NotEqual(t, "", user.ID)
	assert.Equal(t, "ydog", user.UniqueName)
	assert.Equal(t, "yDog", user.DisplayName)
	assert.Equal(t, "ydog@example.com", user.Email)
}

func Test_ユーザーの情報を更新する(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	db, terminate := container.MysqlContainer(
		t,
		ctx,
		container.MySQLcontainerInput(),
	)
	defer terminate()

	repo := repository.NewUserRepository(db)
	s := service.NewUserService(repo)
	err := s.DeleteUser(ctx, "ydog")
	require.Nil(t, err)
	err = s.DeleteUser(ctx, "modified")
	require.Nil(t, err)
	input := model.CreateUserInput{
		UniqueName:  "ydog",
		DisplayName: "yDog",
		Email:       "ydog@example.com",
	}
	id, err := s.CreateUser(ctx, &input)
	require.Nil(t, err)
	require.NotEqual(t, "", id)

	err = s.UpdateUser(ctx, id, &model.UpdateUserInput{
		UniqueName: pkg.PtrStr("modified"),
		Email:      pkg.PtrStr("modified.ydog@example.com"),
	})
	require.Nil(t, err)

	user, err := s.GetUser(ctx, "modified")
	assert.Nil(t, err)
	assert.NotNil(t, user.ID)
	assert.Equal(t, "modified", user.UniqueName)
	assert.Equal(t, "yDog", user.DisplayName)
	assert.Equal(t, "modified.ydog@example.com", user.Email)
}
