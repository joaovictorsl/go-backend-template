package user_test

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joaovictorsl/go-backend-template/internal/config"
	"github.com/joaovictorsl/go-backend-template/internal/database"
	"github.com/joaovictorsl/go-backend-template/internal/entity"
	"github.com/joaovictorsl/go-backend-template/internal/testhelpers"
	"github.com/joaovictorsl/go-backend-template/internal/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RepositoryTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	db          *pgxpool.Pool
	repository  user.Repository
	testUser    entity.User
	ctx         context.Context
}

func (suite *RepositoryTestSuite) SetupSuite() {
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

func (suite *RepositoryTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
}

func (suite *RepositoryTestSuite) SetupTest() {
	_, err := suite.db.Exec(suite.ctx, "DELETE FROM users;")
	if err != nil {
		log.Fatalf("error deleting all users on postgres: %s", err)
	}

	suite.testUser = entity.User{
		ProviderId: "providerId",
		Email:      "email",
	}
}

func (suite *RepositoryTestSuite) TestCreateUserSuccessfully() {
	t := suite.T()
	// Action
	userId, err := suite.repository.CreateUser(suite.ctx, suite.testUser)
	// Assert
	assert.NoError(t, err)
	_, err = uuid.Parse(userId)
	assert.NoError(t, err)
}

func (suite *RepositoryTestSuite) TestCreateUserViolateUniqueConstraint() {
	t := suite.T()
	// Setup
	userId, err := suite.repository.CreateUser(suite.ctx, suite.testUser)
	require.NoError(t, err)
	_, err = uuid.Parse(userId)
	require.NoError(t, err)
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
			assert.Empty(t, userId)
		})
	}
}

func (suite *RepositoryTestSuite) TestGetUserByIdSuccessfully() {
	t := suite.T()
	// Setup
	userId, err := suite.repository.CreateUser(suite.ctx, suite.testUser)
	require.NoError(t, err)
	_, err = uuid.Parse(userId)
	require.NoError(t, err)
	suite.testUser.Id = userId
	// Action
	u, err := suite.repository.GetUserById(suite.ctx, userId)
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, suite.testUser, u)
}

func (suite *RepositoryTestSuite) TestGetUserByIdNotFound() {
	t := suite.T()
	// Action
	u, err := suite.repository.GetUserById(suite.ctx, uuid.NewString())
	// Assert
	assert.ErrorIs(t, err, sql.ErrNoRows)
	assert.Equal(t, entity.User{}, u)
}

func (suite *RepositoryTestSuite) TestGetUserByProviderIdSuccessfully() {
	t := suite.T()
	// Setup
	userId, err := suite.repository.CreateUser(suite.ctx, suite.testUser)
	require.NoError(t, err)
	_, err = uuid.Parse(userId)
	require.NoError(t, err)
	suite.testUser.Id = userId
	// Action
	u, err := suite.repository.GetUserByProviderId(suite.ctx, suite.testUser.ProviderId)
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, suite.testUser, u)
}

func (suite *RepositoryTestSuite) TestGetUserByProviderIdNotFound() {
	t := suite.T()
	// Action
	u, err := suite.repository.GetUserByProviderId(suite.ctx, "fakeId")
	// Assert
	assert.ErrorIs(t, err, sql.ErrNoRows)
	assert.Equal(t, entity.User{}, u)
}

func (suite *RepositoryTestSuite) TestDeleteUserByIdWhenUserExists() {
	t := suite.T()
	// Setup
	userId, err := suite.repository.CreateUser(suite.ctx, suite.testUser)
	require.NoError(t, err)
	_, err = uuid.Parse(userId)
	require.NoError(t, err)
	suite.testUser.Id = userId
	// Action
	err = suite.repository.DeleteUserById(suite.ctx, suite.testUser.Id)
	// Assert
	assert.NoError(t, err)
}

func (suite *RepositoryTestSuite) TestDeleteUserByIdWhenUserDoesNotExist() {
	t := suite.T()
	// Action
	err := suite.repository.DeleteUserById(suite.ctx, uuid.NewString())
	// Assert
	assert.NoError(t, err)
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
