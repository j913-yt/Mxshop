package forms

type PassWordLoginForm struct {
	Mobile    string `from:"mobile" json:"mobile" binding:"required,mobile"` //手机号码格式有规范可寻找，自定义validate
	PassWord  string `form:"password" json:"password" binding:"required,min=3,max=20"`
	Captcha   string `form:"captcha" json:"captcha" binding:"required,min=5,max=5"` //图像验证码
	CaptchaId string `form:"captcha_id" json:"captcha_id" binding:"required"`
}
