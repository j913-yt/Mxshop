package initialize

import (
	"fmt"
	"mxshop_srvs/Mxshop_api/user_web/global"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}
func InitConfig() {
	debug := GetEnvInfo("MXSHOP_DEBUG")

	configFilePerfix := "config"
	configFileName := fmt.Sprintf("Mxshop_api/user_web/%s-pro.yaml", configFilePerfix)
	if debug {
		configFileName = fmt.Sprintf("Mxshop_api/user_web/%s-debug.yaml", configFilePerfix)
	}
	v := viper.New()
	//文件的路径如何设置
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	//这个对象如何在其他文件中使用-全局变量

	if err := v.Unmarshal(&global.ServerConfig); err != nil {
		panic(err)
	}
	fmt.Println(&global.ServerConfig)
	zap.S().Infof("配置信息：%v", global.ServerConfig)
	fmt.Printf("%v", v.Get("name"))

	//viper的功能-动态监控变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Infof("配置文件产生变化：%v", e.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(global.ServerConfig)
		zap.S().Infof("配置信息%v", global.ServerConfig)
	})
}
