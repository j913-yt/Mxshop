package handler

import (
	"context"
	"mxshop_srvs/user_srv/global"
	"mxshop_srvs/user_srv/model"
	"mxshop_srvs/user_srv/proto"

	"gorm.io/gorm"
)

type UserServer struct {
}

func ModelToResponse(user model.User) proto.UserInfoResponse{
	userInfoRsp := proto.UserInfoResponse{
		Id: user.ID,
		PassWord: user.Password,
		NickName:user.NickName ,
		Gender: user.Gender,
		Role: int32(user.Role),

	}
	if user.Birthday != nil{
		userInfoRsp.BrithDay = uint64(user.Birthday.Unix())
	}
	return userInfoRsp
}

func Paginate(page,pageSize int) func(db *gorm.DB) *gorm.DB{
	return func (db *gorm.DB) *gorm.DB {
		if page == 0{
			page = 1
		}

		switch 
		 {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
			
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func (s *UserServer) GetUserList(context.Context, *proto.PageInfo) (*proto.UserListResponse, error) {

		//获取用户列表
		var users []model.User
		result := global.DB.Find(&users)
		if result.Error != nil{
			return nil,result.Error
		}

		rsp := &proto.UserListResponse{
		rsp.Total = int32(result.RowsAffected)

		global.DB.Scopes(Paginate(req.Pn,int(req.PSize))).Find(&users)

		for _,user := range users{
			userInfoRsp := ModelToResponse(user)
			rsp.Data = append(rsp.Data,&userInfoRsp)
		}

		return rsp,nil
	
		}




}
