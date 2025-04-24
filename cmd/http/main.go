package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joaovictorsl/go-backend-template/internal/config"
	"github.com/joaovictorsl/go-backend-template/internal/core/user"
	"github.com/joaovictorsl/go-backend-template/internal/database"
	"github.com/joaovictorsl/go-backend-template/internal/http/handlers"
	"github.com/joaovictorsl/go-backend-template/internal/http/jwt"
)

func main() {
	cfg := config.NewConfig()
	db := database.NewDatabase(cfg)

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	jwtRepository := jwt.NewRepository(db)
	jwtService := jwt.NewService(cfg, jwtRepository)
	authHandler := handlers.NewAuthHandler(cfg, userService, jwtService)

	r := gin.Default()

	r.GET("/auth/oauth/:provider", authHandler.OAuthProvider)
	r.GET("/auth/refresh", authHandler.Refresh)

	if err := r.Run(fmt.Sprintf(":%d", cfg.PORT)); err != nil {
		panic(err)
	}
}
