package api

import (
	"context"
	"fmt"
	"mxshop_srvs/Mxshop_api/user_web/global/reponse"
	"mxshop_srvs/user_srv/proto"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
	ip := "127.0.0.1"
	port := 50051
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList]连接用户服务失败", "msg", err.Error())
	}
	//生成grpc的client并调用接口
	userSrcClient := proto.NewUserClient(userConn)

	rsp, err := userSrcClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    0,
		PSize: 0,
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
