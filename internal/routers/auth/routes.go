package auth

import (
	"gin-realword-example/internal/modules/utils/gin/middleware"
	"gin-realword-example/internal/routers/auth/handlers"

	"github.com/gin-gonic/gin"
	gwm_gin "github.com/slhmy/go-webmods/modules/gin"
)

func RegisterAuthRoutes(router gin.IRouter) {
	routerGroup := router.Group("/auth")
	routerGroup.Use(middleware.RequestID, middleware.Logger, gwm_gin.ErrorHandler)

	routerGroup.
		GET("/github/login", handlers.GithubLogin).
		GET("/github/callback", handlers.GithubCallback)
}
