package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"mxshop_srvs/user_srv/global"
	"mxshop_srvs/user_srv/model"
	"mxshop_srvs/user_srv/proto"
	"strings"
	"time"

	"github.com/anaskhan96/go-password-encoder"
	"github.com/cockroachdb/errors/grpc/status"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
)

type UserServer struct {
	proto.UnimplementedUserServer
}

func ModelToResponse(user model.User) proto.UserInfoResponse {
	userInfoRsp := proto.UserInfoResponse{
		Id:       user.ID,
		PassWord: user.Password,
		NickName: user.NickName,
		Gender:   user.Gender,
		Role:     int32(user.Role),
		Mobile:   user.Mobile,
	}
	if user.Birthday != nil {
		birthDay := uint64(user.Birthday.Unix())
		userInfoRsp.BirthDay = &birthDay
	}
	return userInfoRsp
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10

		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func (s *UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	// Count and page independently so the response data is loaded only once.
	var users []model.User
	var total int64
	result := global.DB.Model(&model.User{}).Count(&total)
	if result.Error != nil {
		return nil, result.Error
	}

	result = global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	rsp := &proto.UserListResponse{Total: int32(total)}

	for _, user := range users {
		userInfoRsp := ModelToResponse(user)
		rsp.Data = append(rsp.Data, &userInfoRsp)
	}

	return rsp, nil

}

// 新建用户方法
func (s *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	//新建用户
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")

	}

	user.Mobile = req.Mobile
	user.NickName = req.NickName

	//密码加密
	options := &password.Options{
		16, 100, 32, sha512.New,
	}
	salt, encodePwd := password.Encode(req.Password, options)
	user.Password = fmt.Sprintf("pbkdf2-sha512$%s$%s", salt, encodePwd)

	result = global.DB.Create(&user)

	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())

	}
	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp, nil

}

// 更新用户
func (s *UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserInfo) (*empty.Empty, error) {
	//个人中心更新用户
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")

	}

	user.NickName = req.NickName
	user.Gender = req.Gender
	if req.BirthDay != nil {
		birthDay := time.Unix(int64(*req.BirthDay), 0)
		user.Birthday = &birthDay
	}

	result = global.DB.Save(user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return &empty.Empty{}, nil

}

// 检查密码
func (s *UserServer) CheckPassWord(ctx context.Context, req *proto.PasswordCheckInfo) (*proto.CheckResponse, error) {
	//校验密码
	options := &password.Options{
		16, 100, 32, sha512.New,
	}
	passwordInfo := strings.Split(req.EncryptedPassword, "$")
	check := password.Verify(req.Password, passwordInfo[1], passwordInfo[2], options)
	return &proto.CheckResponse{
		Success: check,
	}, nil

}
