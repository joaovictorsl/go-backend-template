package repository_test

import (
	"context"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joaovictorsl/go-backend-template/internal/config"
	"github.com/joaovictorsl/go-backend-template/internal/core/entity"
	"github.com/joaovictorsl/go-backend-template/internal/core/errs"
	"github.com/joaovictorsl/go-backend-template/internal/core/user/repository"
	"github.com/joaovictorsl/go-backend-template/internal/database"
	"github.com/joaovictorsl/go-backend-template/internal/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RepositoryTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	db          *pgxpool.Pool
	repository  repository.Repository
	testUser    entity.User
	ctx         context.Context
}

func (suite *RepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	suite.pgContainer = testhelpers.CreatePostgresContainer(suite.ctx)
	cfg := &config.Config{DATABASE_URL: suite.pgContainer.ConnectionString}
	suite.db = database.NewDatabase(cfg)
	suite.repository = repository.New(suite.db)
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
	userId, err := suite.repository.InsertUser(suite.ctx, suite.testUser)
	// Assert
	assert.NoError(t, err)
	_, err = uuid.Parse(userId)
	assert.NoError(t, err)
}

func (suite *RepositoryTestSuite) TestViolatesUniqueConstraintWhenCreatingUsersWithDuplicatedUniqueFields() {
	t := suite.T()
	// Setup
	userId, err := suite.repository.InsertUser(suite.ctx, suite.testUser)
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
			userId, err := suite.repository.InsertUser(suite.ctx, input)
			// Assert
			assert.ErrorIs(t, err, errs.ErrDuplicated)
			assert.Empty(t, userId)
		})
	}
}

func (suite *RepositoryTestSuite) TestSelectUserByIdSuccessfully() {
	t := suite.T()
	// Setup
	userId, err := suite.repository.InsertUser(suite.ctx, suite.testUser)
	require.NoError(t, err)
	_, err = uuid.Parse(userId)
	require.NoError(t, err)
	suite.testUser.Id = userId
	// Action
	u, err := suite.repository.SelectUserById(suite.ctx, userId)
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, suite.testUser, u)
}

func (suite *RepositoryTestSuite) TestSelectUserByIdNotFound() {
	t := suite.T()
	// Action
	u, err := suite.repository.SelectUserById(suite.ctx, uuid.NewString())
	// Assert
	assert.ErrorIs(t, err, errs.ErrNotFound)
	assert.Equal(t, entity.User{}, u)
}

func (suite *RepositoryTestSuite) TestSelectUserByProviderIdSuccessfully() {
	t := suite.T()
	// Setup
	userId, err := suite.repository.InsertUser(suite.ctx, suite.testUser)
	require.NoError(t, err)
	_, err = uuid.Parse(userId)
	require.NoError(t, err)
	suite.testUser.Id = userId
	// Action
	u, err := suite.repository.SelectUserByProviderId(suite.ctx, suite.testUser.ProviderId)
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, suite.testUser, u)
}

func (suite *RepositoryTestSuite) TestSelectUserByProviderIdNotFound() {
	t := suite.T()
	// Action
	u, err := suite.repository.SelectUserByProviderId(suite.ctx, "fakeId")
	// Assert
	assert.ErrorIs(t, err, errs.ErrNotFound)
	assert.Equal(t, entity.User{}, u)
}

func (suite *RepositoryTestSuite) TestDeleteUserByIdWhenUserExists() {
	t := suite.T()
	// Setup
	userId, err := suite.repository.InsertUser(suite.ctx, suite.testUser)
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
