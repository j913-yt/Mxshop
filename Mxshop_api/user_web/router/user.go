package router

import (
	"mxshop_srvs/Mxshop_api/user_web/api"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user")
	zap.S().Info("配置用户相关的url")
	{
		UserRouter.GET("list", api.GetUserList)
	}
}
