package testhelpers

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresContainer struct {
	*postgres.PostgresContainer
	ConnectionString string
}

func CreatePostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	sqlFiles, err := getUpMigrationFiles()
	if err != nil {
		return nil, err
	}

	pgContainer, err := postgres.Run(
		ctx,
		"postgres:17",
		postgres.WithInitScripts(sqlFiles...),
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	return &PostgresContainer{
		PostgresContainer: pgContainer,
		ConnectionString:  connStr,
	}, nil
}

func getUpMigrationFiles() ([]string, error) {
	migrationsPath, err := filepath.Abs(filepath.Join("..", "..", "..", "migrations"))
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(migrationsPath)
	if err != nil {
		return nil, err
	}

	sqlFiles := make([]string, 0)
	for _, a := range entries {
		if strings.Contains(a.Name(), "up") {
			sqlFiles = append(sqlFiles, filepath.Join(migrationsPath, a.Name()))
		}
	}

	return sqlFiles, nil
}
