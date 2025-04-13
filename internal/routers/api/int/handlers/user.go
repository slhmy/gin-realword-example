package handlers

import (
	"net/http"

	gin_utils "gin-realword-example/internal/modules/utils/gin"
	"gin-realword-example/internal/store"

	"github.com/gin-gonic/gin"
	gwm_app "github.com/slhmy/go-webmods/app"
)

func GetUser(ginCtx *gin.Context) {
	id := gwm_app.ID(ginCtx.Param("id"))
	user, err := store.GetUserByID(ginCtx.Request.Context(), id)
	if err != nil {
		gin_utils.AbortWithError(ginCtx, err)
		return
	}
	ginCtx.Negotiate(http.StatusOK, gin.Negotiate{
		Offered: gin_utils.DefaultNegotiateOffered,
		Data:    user,
	})
}
