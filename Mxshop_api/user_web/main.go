package main

import (
	"fmt"
	"mxshop_srvs/Mxshop_api/user_web/initialize"

	"go.uber.org/zap"
)

func main() {
	port := 8021
	//1.初始化logger
	initialize.InitLogger()
	//2.初始化 routers
	Router := initialize.Routers()

	zap.S().Debugf("启动服务器，端口：%d", port)
	err := Router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		zap.S().Panic("启动失败", err.Error())
	}

}
