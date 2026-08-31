package main

import (
	"fmt"
	"mxshop_srvs/Mxshop_api/user_web/initialize"

	"go.uber.org/zap"
)

func main() {
	//初始化 routers
	Router := initialize.Routers()
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	port := 8021
	err := Router.Run(fmt.Sprintf(":%d", port))
	if err != nil {

	}

}
