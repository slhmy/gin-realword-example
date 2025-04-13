package v1

import (
	"gin-realword-example/internal/modules/utils/gin/middleware"
	"gin-realword-example/internal/routers/api/int/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterAPIInternalRoutes(router gin.IRouter) {
	routerGroup := router.Group("/api/internal")
	routerGroup.Use(middleware.RequestID, middleware.Logger, middleware.ErrorHandler)

	routerGroup.Group("/user").
		GET("/:id", handlers.GetUser)
}
