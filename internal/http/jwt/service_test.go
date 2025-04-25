package jwt_test

import (
	"context"
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

type JwtServiceTestSuite struct {
	suite.Suite
	pgContainer   *testhelpers.PostgresContainer
	db            *pgxpool.Pool
	userService   user.Service
	jwtRepository jwt.Repository
	jwtService    jwt.Service
	cfg           *config.Config
	ctx           context.Context
}

func (suite *JwtServiceTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := testhelpers.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	suite.pgContainer = pgContainer

	suite.cfg = &config.Config{
		DATABASE_URL:      pgContainer.ConnectionString,
		JWT_SECRET:        "secret",
		JWT_ISS:           "iss",
		ACCESS_TOKEN_EXP:  1 * time.Second,
		REFRESH_TOKEN_EXP: 2 * time.Second,
	}
	suite.db = database.NewDatabase(suite.cfg)

	jwtRepository := jwt.NewRepository(suite.db)
	if err != nil {
		log.Fatal(err)
	}
	suite.jwtRepository = jwtRepository
	suite.jwtService = jwt.NewService(suite.cfg, jwtRepository)

	userRepository := user.NewRepository(suite.db)
	suite.userService = user.NewService(userRepository)
}

func (suite *JwtServiceTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
}

func (suite *JwtServiceTestSuite) SetupTest() {
	_, err := suite.db.Exec(suite.ctx, "DELETE FROM refresh_tokens;")
	if err != nil {
		log.Fatalf("error deleteing all users on postgres: %s", err)
	}

	_, err = suite.db.Exec(suite.ctx, "DELETE FROM users;")
	if err != nil {
		log.Fatalf("error deleteing all users on postgres: %s", err)
	}
}

func (suite *JwtServiceTestSuite) TestAccessToken() {
	t := suite.T()
	// Setup
	expectedUserId := uint(1)
	// Action
	tok, err := suite.jwtService.NewAccessToken(expectedUserId)
	require.NoError(t, err)
	// Assert
	assert.Equal(t, expectedUserId, tok.UserId)
	assert.Equal(t, suite.cfg.ACCESS_TOKEN_EXP, tok.ExpiresAt.Sub(tok.CreatedAt))
	claims, err := suite.jwtService.ValidateAccessToken(tok.JWT)
	require.NoError(t, err)
	assert.Equal(t, expectedUserId, claims.UserId)
	assert.Equal(t, suite.cfg.JWT_ISS, claims.Issuer)
	// JWT numeric date precision is in seconds
	assert.WithinDuration(t, tok.CreatedAt, claims.IssuedAt.Time, time.Second)
	assert.WithinDuration(t, tok.CreatedAt, claims.NotBefore.Time, time.Second)
	assert.WithinDuration(t, tok.ExpiresAt, claims.ExpiresAt.Time, time.Second)
}

func (suite *JwtServiceTestSuite) TestExpiredAccessToken() {
	t := suite.T()
	// Setup
	expectedUserId := uint(1)

	tok, err := suite.jwtService.NewAccessToken(expectedUserId)
	require.NoError(t, err)
	time.Sleep(suite.cfg.ACCESS_TOKEN_EXP + time.Second)
	// Action
	claims, err := suite.jwtService.ValidateAccessToken(tok.JWT)
	// Assert
	assert.IsType(t, err, jwt.ErrInvalidToken{})
	assert.ErrorContains(t, err, "token is expired")
	assert.Nil(t, claims)
}

func (suite *JwtServiceTestSuite) TestRefreshToken() {
	t := suite.T()
	// Setup
	expectedUserId, err := suite.userService.CreateUser(suite.ctx, entity.User{ProviderId: "providerId", Email: "email"})
	require.NoError(t, err)
	// Action
	tok, err := suite.jwtService.NewRefreshToken(expectedUserId)
	require.NoError(t, err)
	err = suite.jwtService.StoreRefreshToken(suite.ctx, tok)
	require.NoError(t, err)
	// Assert
	assert.Equal(t, expectedUserId, tok.UserId)
	assert.Equal(t, suite.cfg.REFRESH_TOKEN_EXP, tok.ExpiresAt.Sub(tok.CreatedAt))
	claims, err := suite.jwtService.ValidateRefreshToken(suite.ctx, tok.JWT)
	require.NoError(t, err)
	assert.Equal(t, expectedUserId, claims.UserId)
	assert.Equal(t, suite.cfg.JWT_ISS, claims.Issuer)
	// JWT numeric date precision is in seconds
	assert.WithinDuration(t, tok.CreatedAt, claims.IssuedAt.Time, time.Second)
	assert.WithinDuration(t, tok.CreatedAt, claims.NotBefore.Time, time.Second)
	assert.WithinDuration(t, tok.ExpiresAt, claims.ExpiresAt.Time, time.Second)
}

func (suite *JwtServiceTestSuite) TestExpiredRefreshToken() {
	t := suite.T()
	// Setup
	expectedUserId, err := suite.userService.CreateUser(suite.ctx, entity.User{ProviderId: "providerId", Email: "email"})
	require.NoError(t, err)

	tok, err := suite.jwtService.NewRefreshToken(expectedUserId)
	require.NoError(t, err)
	err = suite.jwtService.StoreRefreshToken(suite.ctx, tok)
	require.NoError(t, err)
	time.Sleep(suite.cfg.REFRESH_TOKEN_EXP + time.Second)
	// Action
	claims, err := suite.jwtService.ValidateRefreshToken(suite.ctx, tok.JWT)
	// Assert
	assert.IsType(t, err, jwt.ErrInvalidToken{})
	assert.ErrorContains(t, err, "token is expired")
	assert.Nil(t, claims)
}

func TestJwtServiceTestSuite(t *testing.T) {
	suite.Run(t, new(JwtServiceTestSuite))
}
