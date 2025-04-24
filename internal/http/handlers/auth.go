package handlers

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joaovictorsl/go-backend-template/internal/config"
	"github.com/joaovictorsl/go-backend-template/internal/core/entity"
	"github.com/joaovictorsl/go-backend-template/internal/core/user"
	"github.com/joaovictorsl/go-backend-template/internal/http/jwt"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

type AuthHandler struct {
	cfg         *config.Config
	userService user.Service
	jwtService  jwt.Service
}

func NewAuthHandler(cfg *config.Config, userService user.Service, jwtService jwt.Service) *AuthHandler {
	goth.UseProviders(
		google.New(
			cfg.GOOGLE_CLIENT_ID,
			cfg.GOOGLE_CLIENT_SECRET,
			cfg.GOOGLE_CLIENT_REDIRECT_URI,
		),
	)
	return &AuthHandler{
		cfg:         cfg,
		userService: userService,
		jwtService:  jwtService,
	}
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	cookie, err := c.Request.Cookie("refresh-token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no refresh token provided"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.TIMEOUT)
	defer cancel()

	claims, err := h.jwtService.ValidateRefreshToken(ctx, cookie.Value)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token"})
		return
	}

	newAccessToken, err := h.jwtService.NewToken(claims.UserId, h.cfg.JWT_ISS, h.cfg.JWT_SECRET, h.cfg.ACCESS_TOKEN_EXP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create access token"})
		return
	}

	newRefreshToken, err := h.jwtService.NewToken(claims.UserId, h.cfg.JWT_ISS, h.cfg.JWT_SECRET, h.cfg.REFRESH_TOKEN_EXP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to rotate refresh token"})
		return
	}
	if err := h.jwtService.StoreRefreshToken(ctx, newRefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store new refresh token"})
		return
	}

	setCookie(c.Writer, "token", newAccessToken.JWT, newAccessToken.ExpiresAt, "/", "", false, true)
	setCookie(c.Writer, "refresh-token", newRefreshToken.JWT, newRefreshToken.ExpiresAt, "/", "", false, true)
	c.Status(http.StatusOK)
}

func (h *AuthHandler) OAuthProvider(c *gin.Context) {
	gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		gothic.BeginAuthHandler(c.Writer, c.Request)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.TIMEOUT)
	defer cancel()

	var (
		user entity.User = entity.User{ProviderId: gothUser.UserID, Email: gothUser.Email}
	)
	u, err := h.userService.GetUserByProviderId(ctx, gothUser.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		userId, err := h.userService.CreateUser(ctx, user)
		if err != nil {
			log.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}
		user.Id = userId
	} else if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	} else {
		user = u
	}

	accessToken, err := h.jwtService.NewToken(user.Id, h.cfg.JWT_ISS, h.cfg.JWT_SECRET, h.cfg.ACCESS_TOKEN_EXP)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	refreshToken, err := h.jwtService.NewToken(user.Id, h.cfg.JWT_ISS, h.cfg.JWT_SECRET, h.cfg.REFRESH_TOKEN_EXP)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	err = h.jwtService.StoreRefreshToken(ctx, refreshToken)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	setCookie(c.Writer, "token", accessToken.JWT, accessToken.ExpiresAt, "/", "", false, true)
	setCookie(c.Writer, "refresh-token", refreshToken.JWT, refreshToken.ExpiresAt, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func setCookie(
	w http.ResponseWriter,
	name string,
	value string,
	expiresAt time.Time,
	path string,
	domain string,
	secure bool,
	httpOnly bool,
) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		Expires:  expiresAt,
		Path:     path,
		Domain:   domain,
		SameSite: http.SameSiteLaxMode,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}
