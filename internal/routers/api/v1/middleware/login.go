package middleware

import (
	"errors"
	"net/http"
	"time"

	gin_utils "gin-realword-example/internal/modules/utils/gin"
	"gin-realword-example/internal/routers/shared"
	"gin-realword-example/internal/store"

	"github.com/gin-gonic/gin"
	gwm_app "github.com/slhmy/go-webmods/app"
)

func LoadLoginSession(ginCtx *gin.Context) {
	sessionID, err := ginCtx.Cookie(shared.CookieLoginSessionID)
	if err != nil {
		ginCtx.Next()
		return
	}
	userID, expireAt, err := store.GetUserIDFromLoginSession(ginCtx, sessionID)
	if err != nil {
		ginCtx.SetCookie(shared.CookieLoginSessionID, "", -1, "/", "", false, true)
		ginCtx.Next()
		return
	}
	ginCtx.SetCookie(shared.CookieLoginSessionID, sessionID, int(time.Until(*expireAt).Seconds()), "/", "", false, true)
	ginCtx.Set(string(shared.ContextKeyUserID), *userID)
	ginCtx.Next()
}

func RequireLoginSession(ginCtx *gin.Context) {
	userID := ginCtx.GetUint(string(shared.ContextKeyUserID))
	if userID == 0 {
		gin_utils.AbortWithError(ginCtx, gwm_app.ServiceError{
			Err:            errors.New("unauthorized"),
			HttpStatusCode: http.StatusUnauthorized,
		})
		return
	}
	ginCtx.Next()
}
