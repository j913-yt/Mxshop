package forms

type PassWordLoginForm struct {
	Mobile   string `from:"mobile" json:"mobile" binding:"required,mobile"` //手机号码格式有规范可寻找，自定义validate
	PassWord string `form:"password" json:"password" binding:"required,min=3,max=20"`
}
