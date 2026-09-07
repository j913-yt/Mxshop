package initialize

import (
	"fmt"
	"mxshop_srvs/Mxshop_api/user_web/global"
	"reflect"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

// 第一步编写一个初始化翻译器函数
func InitTrans() error {
	//初始化中文环境
	zhT := zh.New()
	//配置一个通用的翻译器(备用语言是中文，当前支持的语言也是中文）
	uni := ut.New(zhT, zhT)
	//将全局翻译器设为中文
	var ok bool
	global.Trans, ok = uni.GetTranslator("zh")

	if !ok {
		return fmt.Errorf("获取翻译器失败")

	}

	//4.从gin中提取出底层的validator引擎
	//注意这里使用了类型断言，将Gin底层的验证器转换为validator.Validate实例
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {

		//核心改造一，注册一个获取自定义字段名的函数
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			//去结构体里找名叫“label”的标签
			name := fld.Tag.Get("label")
			//如果找到了，就用这个中文名代替原本的英文字段名
			if name != "" {
				return name
			}
			//如果没写label，，就按原样返回
			return fld.Name
		})
		err := zhTranslations.RegisterDefaultTranslations(v, global.Trans)
		if err != nil {
			return err
		}
	}
	return nil

}
