package main

import (
	"fmt"
	"mxshop_srvs/Mxshop_api/user_web/global"
	"mxshop_srvs/Mxshop_api/user_web/initialize"

	"go.uber.org/zap"
)

func main() {

	//1.初始化logger
	initialize.InitLogger()
	//2.初始化配置文件
	initialize.InitConfig()
	//3.初始化 routers
	Router := initialize.Routers()

	zap.S().Debugf("启动服务器，端口：%d", global.ServerConfig.Port)
	err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port))
	if err != nil {
		zap.S().Panic("启动失败", err.Error())
	}

}
