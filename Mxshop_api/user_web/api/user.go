package api

import (
	"context"
	"fmt"
	"mxshop_srvs/Mxshop_api/user_web/forms"
	"mxshop_srvs/Mxshop_api/user_web/global"
	"mxshop_srvs/Mxshop_api/user_web/global/reponse"
	"mxshop_srvs/Mxshop_api/user_web/middlewares"
	"mxshop_srvs/Mxshop_api/user_web/models"
	"mxshop_srvs/user_srv/proto"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	//将grpc的code转换成http的状态码
	if err != nil {
		e, ok := status.FromError(err)
		if ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": e.Code(),
				})

				return

			}
		}
	}
}

func GetUserList(ctx *gin.Context) {

	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList]连接用户服务失败", "msg", err.Error())
	}

	claims, _ := ctx.Get("claims")
	currentUser := claims.(*models.CustomClaims)
	zap.S().Infof("访问用户:%d", currentUser.ID)

	//生成grpc的client并调用接口
	userSrcClient := proto.NewUserClient(userConn)

	pn := ctx.DefaultQuery("pn", "0")
	pSize := ctx.DefaultQuery("psize", "10")
	pnInt, _ := strconv.Atoi(pn)
	pSizeInt, _ := strconv.Atoi(pSize)
	rsp, err := userSrcClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
	})
	if err != nil {
		zap.S().Errorw("GetUserList 查询用户列表失败")
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		//data := make(map[string]interface{})

		user := reponse.UserResponse{
			Id:       value.Id,
			NickName: value.NickName,
			Gender:   value.Gender,
			Mobile:   value.Mobile,
		}
		if value.BirthDay != nil {
			birthday := reponse.JsonTime(time.Unix(int64(*value.BirthDay), 0))
			user.Birthday = &birthday
		}
		//data["id"] = value.Id
		//data["name"] = value.NickName
		//data["birthday"] = value.BirthDay
		//data["gender"] = value.Gender
		//data["mobile"] = value.Mobile
		result = append(result, user)
	}
	ctx.JSON(http.StatusOK, result)

}

func PassWordLogin(c *gin.Context) {
	//表单验证
	PassWordLoginForm := forms.PassWordLoginForm{}
	err := c.ShouldBindJSON(&PassWordLoginForm)
	//如何返回错误信息
	if err != nil {
		//拦截错误，进行翻译
		//我们需要判断，这个错误是不是validator验证失败抛出的错误
		HandleValidatorError(c, err)
		return
	}
	if !store.Verify(PassWordLoginForm.CaptchaId, PassWordLoginForm.Captcha, true) {
		c.JSON(http.StatusBadRequest, gin.H{
			"Captcha": "验证码错误",
		})
		return
	}
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList]连接用户服务失败", "msg", err.Error())
	}
	//生成grpc的client并调用接口
	userSrcClient := proto.NewUserClient(userConn)

	//登录的逻辑
	rsp, err := userSrcClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: PassWordLoginForm.Mobile,
	})
	if err != nil {
		e, ok := status.FromError(err)
		if ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusBadRequest, map[string]string{
					"mobile": "用户不存在",
				})
			default:
				c.JSON(http.StatusInternalServerError, map[string]string{
					"mobile": "登录失败",
				})

			}
			return
		}
	} else {
		//只是查询到用户而已，并没有检查密码
		passReq, passErr := userSrcClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
			Password:          PassWordLoginForm.PassWord,
			EncryptedPassword: rsp.PassWord,
		})
		if passErr != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{
				"password": "登录失败",
			})
		} else {
			if passReq.Success {
				//生成token
				j := middlewares.NewJWT()
				claims := models.CustomClaims{
					ID:          uint(rsp.Id),
					NickName:    rsp.NickName,
					AuthorityId: uint(rsp.Role),
					StandardClaims: jwt.StandardClaims{
						NotBefore: time.Now().Unix(),               //签名的生效时间
						ExpiresAt: time.Now().Unix() + 60*60*24*30, //30天过期
						Issuer:    "imooc",
					},
				}
				token, err := j.CreateToken(claims)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"msg": "生成token失败",
					})
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"id":         rsp.Id,
					"nick_name":  rsp.NickName,
					"token":      token,
					"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000,
				})
			} else {
				c.JSON(http.StatusBadRequest, map[string]string{
					"msg": "登录失败",
				})
			}

		}

	}

}

func HandleValidatorError(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		//不是validator的错误，说明是json格式本身的问题，直接返回
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
		return
	}
	//核心改造3，手动格式化错误字典，去掉结构体前缀
	c.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(global.Trans)),
	})
}

// removeTopStruct 去掉翻译后错误信息里的结构体前缀
// 例如 "PassWordLoginForm.mobile" 只保留 "mobile"
func removeTopStruct(fields map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fields {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}
