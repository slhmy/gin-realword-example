package middleware

import (
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
	userID := ginCtx.GetString(string(shared.ContextKeyUserID))
	if userID == "" {
		gin_utils.AbortWithError(ginCtx, gwm_app.NewForbiddenError("not logged in"))
		return
	}
	ginCtx.Next()
}
