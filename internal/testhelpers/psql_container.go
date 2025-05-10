package testhelpers

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	once        = &sync.Once{}
	pgContainer *PostgresContainer
)

type PostgresContainer struct {
	*postgres.PostgresContainer
	ConnectionString string
}

func CreatePostgresContainer(ctx context.Context) *PostgresContainer {
	once.Do(func() {
		sqlFiles, err := getUpMigrationFiles()
		if err != nil {
			panic(err)
		}

		container, err := postgres.Run(
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
			panic(err)
		}
		connStr, err := container.ConnectionString(ctx, "sslmode=disable")
		if err != nil {
			panic(err)
		}

		pgContainer = &PostgresContainer{
			PostgresContainer: container,
			ConnectionString:  connStr,
		}
	})

	return pgContainer
}

func getUpMigrationFiles() ([]string, error) {
	migrationsPath, err := filepath.Abs(filepath.Join("..", "..", "..", "..", "migrations"))
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
