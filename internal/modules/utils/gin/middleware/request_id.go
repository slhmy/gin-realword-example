package middleware

import (
	"log/slog"

	"gin-realword-example/internal/modules/core"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	RequestIDHeaderKey  = "X-Request-ID"
	RequestIDContextKey = "requestID"
)

func RequestID(ginCtx *gin.Context) {
	requestID := ginCtx.Request.Header.Get(RequestIDHeaderKey)
	if requestID == "" {
		requestID = uuid.NewString()
	}
	ginCtx.Writer.Header().Set(RequestIDHeaderKey, requestID)
	core.AppendLogFieldToGinCtx(ginCtx, slog.String("request_id", requestID))
	ginCtx.Next()
}
