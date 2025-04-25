package jwt_test

import (
	"context"
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joaovictorsl/go-backend-template/internal/config"
	"github.com/joaovictorsl/go-backend-template/internal/core/entity"
	"github.com/joaovictorsl/go-backend-template/internal/core/user"
	"github.com/joaovictorsl/go-backend-template/internal/database"
	"github.com/joaovictorsl/go-backend-template/internal/http/jwt"
	"github.com/joaovictorsl/go-backend-template/internal/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type JwtRepositoryTestSuite struct {
	suite.Suite
	pgContainer   *testhelpers.PostgresContainer
	db            *pgxpool.Pool
	userService   user.Service
	jwtRepository jwt.Repository
	testToken     jwt.Token
	ctx           context.Context
}

func (suite *JwtRepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := testhelpers.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	suite.pgContainer = pgContainer

	cfg := &config.Config{DATABASE_URL: pgContainer.ConnectionString}
	suite.db = database.NewDatabase(cfg)

	repository := jwt.NewRepository(suite.db)
	if err != nil {
		log.Fatal(err)
	}
	suite.jwtRepository = repository

	userRepository := user.NewRepository(suite.db)
	suite.userService = user.NewService(userRepository)
}

func (suite *JwtRepositoryTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
}

func (suite *JwtRepositoryTestSuite) SetupTest() {
	_, err := suite.db.Exec(suite.ctx, "DELETE FROM refresh_tokens;")
	if err != nil {
		log.Fatalf("error deleteing all users on postgres: %s", err)
	}

	_, err = suite.db.Exec(suite.ctx, "DELETE FROM users;")
	if err != nil {
		log.Fatalf("error deleteing all users on postgres: %s", err)
	}

	createdAt := time.Date(2025, time.January, 20, 0, 0, 0, 0, time.Local)
	expiresAt := createdAt.Add(15 * time.Second)
	suite.testToken = jwt.Token{
		UserId:    1,
		JWT:       "jwt",
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	}
}

func (suite *JwtRepositoryTestSuite) TestStoreTokenNonExistentUser() {
	t := suite.T()
	// Action
	err := suite.jwtRepository.StoreToken(suite.ctx, suite.testToken)
	// Assert
	assert.Error(t, err)
}

func (suite *JwtRepositoryTestSuite) TestStoreTokenExistentUser() {
	t := suite.T()
	// Setup
	id, err := suite.userService.CreateUser(suite.ctx, entity.User{ProviderId: "providerId", Email: "email"})
	require.NoError(t, err)
	suite.testToken.UserId = id
	// Action
	err = suite.jwtRepository.StoreToken(suite.ctx, suite.testToken)
	// Assert
	assert.NoError(t, err)
}

func (suite *JwtRepositoryTestSuite) TestStoreTokenWithCreatedAtAfterExpiresAt() {
	t := suite.T()
	// Setup
	id, err := suite.userService.CreateUser(suite.ctx, entity.User{ProviderId: "providerId", Email: "email"})
	require.NoError(t, err)
	suite.testToken.UserId = id
	createdAt := suite.testToken.CreatedAt
	suite.testToken.CreatedAt = suite.testToken.ExpiresAt
	suite.testToken.ExpiresAt = createdAt
	// Action
	err = suite.jwtRepository.StoreToken(suite.ctx, suite.testToken)
	// Assert
	assert.Error(t, err)
}

func (suite *JwtRepositoryTestSuite) TestGetNonExistentToken() {
	t := suite.T()
	// Action
	tok, err := suite.jwtRepository.GetToken(suite.ctx, 1)
	// Assert
	assert.Error(t, err)
	assert.Equal(t, jwt.Token{}, tok)
}

func (suite *JwtRepositoryTestSuite) TestGetToken() {
	t := suite.T()
	// Setup
	userId, err := suite.userService.CreateUser(suite.ctx, entity.User{ProviderId: "providerId", Email: "email"})
	require.NoError(t, err)
	suite.testToken.UserId = userId
	err = suite.jwtRepository.StoreToken(suite.ctx, suite.testToken)
	require.NoError(t, err)
	// Action
	tok, err := suite.jwtRepository.GetToken(suite.ctx, userId)
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, suite.testToken.UserId, tok.UserId)
	assert.Equal(t, suite.testToken.JWT, tok.JWT)
	assert.WithinDuration(t, suite.testToken.CreatedAt, tok.CreatedAt, time.Microsecond)
	assert.WithinDuration(t, suite.testToken.ExpiresAt, tok.ExpiresAt, time.Microsecond)
}

func (suite *JwtRepositoryTestSuite) TestDeleteTokensWhenUserIsDeleted() {
	t := suite.T()
	// Setup
	userId, err := suite.userService.CreateUser(suite.ctx, entity.User{ProviderId: "providerId", Email: "email"})
	require.NoError(t, err)
	suite.testToken.UserId = userId
	err = suite.jwtRepository.StoreToken(suite.ctx, suite.testToken)
	require.NoError(t, err)
	err = suite.userService.DeleteUserById(suite.ctx, userId)
	require.NoError(t, err)
	// Action
	tok, err := suite.jwtRepository.GetToken(suite.ctx, userId)
	// Assert
	assert.ErrorIs(t, err, sql.ErrNoRows)
	assert.Equal(t, jwt.Token{}, tok)
}

func TestJwtRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(JwtRepositoryTestSuite))
}
