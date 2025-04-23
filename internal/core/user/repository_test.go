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
	testUser    entity.User
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

func (suite *UserRepositoryTestSuite) SetupTest() {
	_, err := suite.db.Exec(suite.ctx, "DELETE FROM users;")
	if err != nil {
		log.Fatalf("error deleteing all users on postgres: %s", err)
	}

	suite.testUser = entity.User{
		ProviderId: "providerId",
		Email:      "email",
	}
}

func (suite *UserRepositoryTestSuite) TestCreateUserViolateUniqueConstraint() {
	t := suite.T()
	// Setup
	userId, err := suite.repository.CreateUser(suite.ctx, suite.testUser)
	require.NoError(t, err)
	require.Greater(t, userId, uint(0))
	testTable := map[string]entity.User{
		"duplicated provider id": {
			ProviderId: suite.testUser.ProviderId,
			Email:      "otherEmail",
		},
		"duplicated email": {
			ProviderId: "otherProviderId",
			Email:      suite.testUser.Email,
		},
	}
	for testCase, input := range testTable {
		t.Run(testCase, func(t *testing.T) {
			// Action
			userId, err := suite.repository.CreateUser(suite.ctx, input)
			// Assert
			assert.Error(t, err)
			assert.EqualValues(t, 0, userId)
		})
	}
}

func (suite *UserRepositoryTestSuite) TestCreateUserSuccessfully() {
	t := suite.T()
	// Action
	userId, err := suite.repository.CreateUser(suite.ctx, suite.testUser)
	// Assert
	assert.NoError(t, err)
	assert.Greater(t, userId, uint(0))
}

func (suite *UserRepositoryTestSuite) TestGetUserByIdNotFound() {
	t := suite.T()
	// Action
	userId, err := suite.repository.GetUserById(suite.ctx, 1)
	// Assert
	assert.ErrorIs(t, err, sql.ErrNoRows)
	assert.Equal(t, entity.User{}, userId)
}

func (suite *UserRepositoryTestSuite) TestGetUserByIdSuccessfully() {
	t := suite.T()
	// Setup
	userId, err := suite.repository.CreateUser(suite.ctx, suite.testUser)
	suite.testUser.Id = userId
	require.NoError(t, err)
	require.Greater(t, userId, uint(0))
	// Action
	user, err := suite.repository.GetUserById(suite.ctx, userId)
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, suite.testUser, user)
}

func (suite *UserRepositoryTestSuite) TestGetUserByProviderIdNotFound() {
	t := suite.T()
	// Action
	user, err := suite.repository.GetUserByProviderId(suite.ctx, "fakeId")
	// Assert
	assert.ErrorIs(t, err, sql.ErrNoRows)
	assert.Equal(t, entity.User{}, user)
}

func (suite *UserRepositoryTestSuite) TestGetUserByProviderIdSuccessfully() {
	t := suite.T()
	// Setup
	userId, err := suite.repository.CreateUser(suite.ctx, suite.testUser)
	suite.testUser.Id = userId
	require.NoError(t, err)
	require.Greater(t, userId, uint(0))
	// Action
	user, err := suite.repository.GetUserByProviderId(suite.ctx, suite.testUser.ProviderId)
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, suite.testUser, user)
}

func TestCustomerRepoTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
