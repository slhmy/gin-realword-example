package main

import (
	"log/slog"

	"gin-realword-example/internal/modules/core"
	"gin-realword-example/internal/modules/shared"
	gin_utils "gin-realword-example/internal/modules/utils/gin"
	v1 "gin-realword-example/internal/routers/api/v1"
	"gin-realword-example/internal/routers/auth"
	"gin-realword-example/internal/routers/website"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func main() {
	if core.Env() != core.EnvDevelopment {
		gin.SetMode(gin.ReleaseMode)
	}
	ginEngine := gin.New()
	binding.Validator = new(gin_utils.DefaultValidator)
	ginEngine.UseH2C = true
	website.LoadHTMLFiles(ginEngine)

	v1.RegisterAPIV1Routes(ginEngine)
	auth.RegisterAuthRoutes(ginEngine)
	website.RegisterWebsiteRoutes(ginEngine)
	port := core.ConfigStore.GetString(shared.ConfigKeyWebPort)

	slog.Info("Listening",
		slog.String("port", port),
		slog.Bool("enable_tls", false),
	)
	err := ginEngine.Run(":" + port)
	if err != nil {
		panic(err)
	}
}

func init() {
	core.ServiceRole = "web"
	slog.Info("Initialized",
		slog.String("env", core.Env()),
		slog.String("projectDir", core.GetConfigDir()),
	)
}
