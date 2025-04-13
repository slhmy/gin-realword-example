package gin_utils

import "github.com/gin-gonic/gin"

func AbortWithError(ginCtx *gin.Context, err error) {
	if err != nil {
		_ = ginCtx.Error(err)
	}
	ginCtx.Abort()
}
