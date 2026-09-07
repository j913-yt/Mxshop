package global

import (
	"mxshop_srvs/Mxshop_api/user_web/config"

	ut "github.com/go-playground/universal-translator"
)

var (
	// 全局翻译变量
	Trans ut.Translator

	ServerConfig *config.ServerConfig = &config.ServerConfig{}
)
