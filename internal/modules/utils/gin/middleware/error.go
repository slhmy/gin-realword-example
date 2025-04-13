package middleware

import (
	"fmt"
	"log/slog"
	"runtime/debug"

	gin_utils "gin-realword-example/internal/modules/utils/gin"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error string `form:"error" json:"error"`
}

func ErrorHandler(ginCtx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			slog.ErrorContext(ginCtx, string(stack),
				slog.Any("error", r),
				slog.Bool("panic", true),
			)
			_ = ginCtx.Error(fmt.Errorf("panic: %v", r))
		}
		if len(ginCtx.Errors) == 0 {
			return
		}
		err := gin_utils.WrapServiceError(ginCtx.Errors.Last().Err)
		slog.ErrorContext(ginCtx, err.Error())
		ginCtx.Negotiate(err.HttpStatus, gin.Negotiate{
			Offered: gin_utils.DefaultNegotiateOffered,
			Data:    ErrorResponse{Error: err.Error()},
		})
	}()
	ginCtx.Next()
}
