package handlers

import (
	"net/http"
	"strconv"

	gin_utils "gin-realword-example/internal/modules/utils/gin"
	"gin-realword-example/internal/store"

	"github.com/gin-gonic/gin"
)

func GetUser(ginCtx *gin.Context) {
	id, err := strconv.Atoi(ginCtx.Param("id"))
	if err != nil {
		gin_utils.AbortWithError(ginCtx, err)
		return
	}
	user, err := store.GetUserByID(ginCtx.Request.Context(), uint(id))
	if err != nil {
		gin_utils.AbortWithError(ginCtx, gin_utils.WrapServiceError(err))
		return
	}
	ginCtx.Negotiate(http.StatusOK, gin.Negotiate{
		Offered: gin_utils.DefaultNegotiateOffered,
		Data:    user,
	})
}
