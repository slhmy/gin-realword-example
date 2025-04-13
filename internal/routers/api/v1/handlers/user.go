package handlers

import (
	"net/http"

	gin_utils "gin-realword-example/internal/modules/utils/gin"
	"gin-realword-example/internal/routers/shared"
	"gin-realword-example/internal/store"

	"github.com/gin-gonic/gin"
)

// GetUserMe
//
//	@Id				GetUserMe
//	@Summary		Get current user
//	@Description	Get current user
//	@Tags			User
//	@Produce		json
//	@Success		200	{object}	models.User
//	@Failure		401	{object}	middleware.ErrorResponse
//	@Failure		500	{object}	middleware.ErrorResponse
//	@Router			/user/me [get]
func GetUserMe(ginCtx *gin.Context) {
	userID := ginCtx.GetUint(string(shared.ContextKeyUserID))
	if userID == 0 {
		gin_utils.AbortWithError(ginCtx, gin_utils.ErrUnauthorized)
	}
	user, err := store.GetUserByID(ginCtx, uint(userID))
	if err != nil {
		gin_utils.AbortWithError(ginCtx, err)
		return
	}
	ginCtx.Negotiate(http.StatusOK, gin.Negotiate{
		Offered: gin_utils.DefaultNegotiateOffered,
		Data:    user,
	})
}
