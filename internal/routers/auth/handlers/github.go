package handlers

import (
	"net/http"
	"time"

	"gin-realword-example/internal/models"
	github_client "gin-realword-example/internal/modules/clients/github"
	gin_utils "gin-realword-example/internal/modules/utils/gin"
	"gin-realword-example/internal/routers/shared"
	"gin-realword-example/internal/store"

	"github.com/gin-gonic/gin"
	gwm_app "github.com/slhmy/go-webmods/app"
)

const (
	githubCallbackUrl = "/auth/github/callback"
)

func GithubLogin(ginCtx *gin.Context) {
	u, err := github_client.BuildOauthEntryUrl(githubCallbackUrl)
	if err != nil {
		gin_utils.AbortWithError(ginCtx, err)
		return
	}
	ginCtx.Redirect(http.StatusFound, u.String())
}

func GithubCallback(ginCtx *gin.Context) {
	code := ginCtx.Query("code")
	if code == "" {
		gin_utils.AbortWithError(ginCtx, gwm_app.NewBadRequestError("missing code query parameter"))
		return
	}
	tokenResponse, err := github_client.GetAccessToken(code)
	if err != nil {
		gin_utils.AbortWithError(ginCtx, err)
		return
	}
	githubUser, err := github_client.GetUser(tokenResponse.AccessToken)
	if err != nil {
		gin_utils.AbortWithError(ginCtx, err)
		return
	}
	id, err := store.UpsertUser(ginCtx, models.CreateUserRequest{
		Username: githubUser.Login,
		Email:    githubUser.Email,
	})
	if err != nil {
		gin_utils.AbortWithError(ginCtx, err)
		return
	}
	if id == nil {
		user, err := store.GetUserByEmail(ginCtx, githubUser.Email)
		if err != nil {
			gin_utils.AbortWithError(ginCtx, err)
			return
		}
		id = &user.ID
	}
	session, expireAt, err := store.GenerateLoginSession(ginCtx, *id)
	if err != nil {
		gin_utils.AbortWithError(ginCtx, err)
		return
	}
	ginCtx.SetCookie(shared.CookieLoginSessionID, *session, int(time.Until(*expireAt).Seconds()), "/", "", false, true)
	ginCtx.Redirect(http.StatusFound, "/")
}
