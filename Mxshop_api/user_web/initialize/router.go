package initialize

import (
	router2 "Mxshop_api/user_web/router"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	ApiGroup := Router.Group("/v1")
	router2.InitUserRouter(ApiGroup)
}
