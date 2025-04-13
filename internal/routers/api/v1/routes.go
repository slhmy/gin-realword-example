package v1

import (
	common_middleware "gin-realword-example/internal/modules/utils/gin/middleware"
	"gin-realword-example/internal/routers/api/v1/handlers"
	"gin-realword-example/internal/routers/api/v1/middleware"

	"github.com/gin-gonic/gin"
	gwm_gin "github.com/slhmy/go-webmods/modules/gin"
)

// @title		Gin Real World Example API
// @version	1.0
// @BasePath	/api/v1
func RegisterAPIV1Routes(router gin.IRouter) {
	routerGroup := router.Group("/api/v1")
	routerGroup.Use(common_middleware.RequestID, common_middleware.Logger, gwm_gin.ErrorHandler)

	routerGroup.Group("/user").
		Use(middleware.LoadLoginSession, middleware.RequireLoginSession).
		GET("/me", handlers.GetUserMe)
	routerGroup.Group("/chat").
		Use(middleware.LoadLoginSession, middleware.RequireLoginSession).
		GET("/stream", handlers.ChatStream)
}
