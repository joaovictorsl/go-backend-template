package user_test

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joaovictorsl/go-backend-template/internal/config"
	"github.com/joaovictorsl/go-backend-template/internal/core/entity"
	"github.com/joaovictorsl/go-backend-template/internal/core/user"
	"github.com/joaovictorsl/go-backend-template/internal/database"
	"github.com/joaovictorsl/go-backend-template/internal/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	db          *pgxpool.Pool
	repository  user.Repository
	ctx         context.Context
}

func (suite *UserRepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := testhelpers.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	suite.pgContainer = pgContainer

	cfg := &config.Config{DATABASE_URL: pgContainer.ConnectionString}
	suite.db = database.NewDatabase(cfg)

	repository := user.NewRepository(suite.db)
	if err != nil {
		log.Fatal(err)
	}
	suite.repository = repository
}

func (suite *UserRepositoryTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
}

func (suite *UserRepositoryTestSuite) TearDownTest() {
	_, err := suite.db.Exec(suite.ctx, "TRUNCATE TABLE users RESTART IDENTITY CASCADE;")
	if err != nil {
		log.Fatalf("error restoring postgres container to snapshot: %s", err)
	}
}

func (suite *UserRepositoryTestSuite) TestCreateUserViolateUniqueConstraint() {
	t := suite.T()
	// Setup
	userId, err := suite.repository.CreateUser(suite.ctx, entity.User{
		GoogleId: "fakeId",
		Email:    "test@test.com",
		Username: "test",
	})
	require.NoError(t, err)
	require.Equal(t, uint(1), userId)
	testTable := map[string]entity.User{
		"duplicated google id": {
			GoogleId: "fakeId",
			Email:    "otherTest@test.com",
			Username: "otherTest",
		},
		"duplicated email": {
			GoogleId: "otherFakeId",
			Email:    "test@test.com",
			Username: "otherTest",
		},
		"duplicated username": {
			GoogleId: "otherFakeId",
			Email:    "otherTest@test.com",
			Username: "test",
		},
	}
	for testCase, input := range testTable {
		t.Run(testCase, func(t *testing.T) {
			// Action
			_, err := suite.repository.CreateUser(suite.ctx, input)
			// Assert
			assert.Error(t, err)
		})
	}
}

func (suite *UserRepositoryTestSuite) TestCreateUserSuccessfully() {
	t := suite.T()
	// Action
	userId, err := suite.repository.CreateUser(suite.ctx, entity.User{
		GoogleId: "fakeId",
		Email:    "test@test.com",
		Username: "test",
	})
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, uint(1), userId)
}

func (suite *UserRepositoryTestSuite) TestGetUserByIdNotFound() {
	t := suite.T()
	// Action
	_, err := suite.repository.GetUserById(suite.ctx, 1)
	// Assert
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

func (suite *UserRepositoryTestSuite) TestGetUserByIdSuccessfully() {
	t := suite.T()
	// Setup
	userId, err := suite.repository.CreateUser(suite.ctx, entity.User{
		GoogleId: "fakeId",
		Email:    "test@test.com",
		Username: "test",
	})
	require.NoError(t, err)
	require.Equal(t, uint(1), userId)
	// Action
	user, err := suite.repository.GetUserById(suite.ctx, 1)
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, uint(1), user.Id)
}

func (suite *UserRepositoryTestSuite) TestGetUserByGoogleIdNotFound() {
	t := suite.T()
	// Action
	_, err := suite.repository.GetUserByGoogleId(suite.ctx, "fakeId")
	// Assert
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

func (suite *UserRepositoryTestSuite) TestGetUserByGoogleIdSuccessfully() {
	t := suite.T()
	// Setup
	userId, err := suite.repository.CreateUser(suite.ctx, entity.User{
		GoogleId: "fakeId",
		Email:    "test@test.com",
		Username: "test",
	})
	require.NoError(t, err)
	require.Equal(t, uint(1), userId)
	// Action
	user, err := suite.repository.GetUserByGoogleId(suite.ctx, "fakeId")
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "fakeId", user.GoogleId)
}

func TestCustomerRepoTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
