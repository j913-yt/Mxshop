package main

import (
	"fmt"
	"mxshop_srvs/Mxshop_api/user_web/global"
	"mxshop_srvs/Mxshop_api/user_web/initialize"
	myvalidator "mxshop_srvs/Mxshop_api/user_web/validator"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func main() {

	//1.初始化logger
	initialize.InitLogger()
	//2.初始化配置文件
	initialize.InitConfig()
	//3.初始化 routers
	Router := initialize.Routers()

	//5.注册验证器
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		v.RegisterValidation("mobile", myvalidator.ValidateMobile)
	}
	v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} must have a value!", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())

		return t
	})
	//4.初始化翻译
	err2 := initialize.InitTrans()

	if err2 != nil {
		panic(err2)
	}
	zap.S().Debugf("启动服务器，端口：%d", global.ServerConfig.Port)

	err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port))
	if err != nil {
		zap.S().Panic("启动失败", err.Error())
	}

}
