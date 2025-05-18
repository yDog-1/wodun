package container

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

type mySQLcontainerInput struct {
	image    string
	dbName   string
	user     string
	password string
}

type MySQLcontainerOption func(*mySQLcontainerInput)

func MySQLcontainerInput(options ...MySQLcontainerOption) *mySQLcontainerInput {
	input := &mySQLcontainerInput{
		image:    "mysql:8.0-debian",
		dbName:   "wodun",
		user:     "testuser",
		password: "password",
	}
	for _, o := range options {
		o(input)
	}
	return input
}

func WithImage(image string) MySQLcontainerOption {
	return func(input *mySQLcontainerInput) {
		input.image = image
	}
}

func WithDBName(dbName string) MySQLcontainerOption {
	return func(input *mySQLcontainerInput) {
		input.dbName = dbName
	}
}

func WithUser(user string) MySQLcontainerOption {
	return func(input *mySQLcontainerInput) {
		input.user = user
	}
}

func WithPassword(password string) MySQLcontainerOption {
	return func(input *mySQLcontainerInput) {
		input.password = password
	}
}

func MysqlContainer(t *testing.T, ctx context.Context, input *mySQLcontainerInput) (db *sql.DB, terminate func()) {
	t.Helper()

	mysqlContainer, err := mysql.Run(ctx,
		input.image,
		mysql.WithDatabase(input.dbName),
		mysql.WithUsername(input.user),
		mysql.WithPassword(input.password),
	)

	require.NoError(t, err)
	connectionString, err := mysqlContainer.ConnectionString(ctx)
	require.NoError(t, err)

	d, err := sql.Open("mysql", connectionString)
	require.NoError(t, err)
	// マイグレーションファイルを取得
	migrationsDir := "../sql/migrations"
	p, err := goose.NewProvider(goose.DialectMySQL, d, os.DirFS(migrationsDir))
	require.NoError(t, err)
	// マイグレーションを実行
	_, err = p.Up(ctx)
	require.NoError(t, err)

	return d, func() {
		if err := testcontainers.TerminateContainer(mysqlContainer); err != nil {
			require.NoError(t, err)
		}
	}
}
