package website

import (
	"gin-realword-example/internal/modules/core"
	"path"

	"github.com/gin-gonic/gin"
)

const websiteDist = "website/dist"

func html(ginCtx *gin.Context) {
	ginCtx.HTML(200, "index.html", nil)
}

func LoadHTMLFiles(ginEngine *gin.Engine) {
	ginEngine.LoadHTMLFiles(path.Join(core.GetProjectDir(), websiteDist, "index.html"))
}

func RegisterWebsiteRoutes(router gin.IRouter) {
	routerGroup := router.Group("")

	routerGroup.
		GET("", html).
		StaticFile("/vite.svg", path.Join(core.GetProjectDir(), websiteDist, "vite.svg")).
		Static("/assets", path.Join(core.GetProjectDir(), websiteDist, "assets"))
}
