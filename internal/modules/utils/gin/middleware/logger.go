package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger(ginCtx *gin.Context) {
	path := ginCtx.Request.URL.Path
	clientIP := ginCtx.ClientIP()
	method := ginCtx.Request.Method
	userAgent := ginCtx.Request.UserAgent()

	start := time.Now()
	ginCtx.Next()
	end := time.Now()

	latency := int(end.Sub(start) / time.Millisecond)
	statusCode := ginCtx.Writer.Status()

	if skipLog(ginCtx, statusCode) {
		return
	}

	slog.InfoContext(ginCtx, "HTTP request",
		slog.Int("status_code", statusCode),
		slog.Int("latency", latency),
		slog.String("user_agent", userAgent),
		slog.String("method", method),
		slog.String("client_ip", clientIP),
		slog.String("path", path),
	)
}

func skipLog(ginCtx *gin.Context, statusCode int) bool {
	return (statusCode == 200 && ginCtx.Request.URL.Path == "/health")
}
