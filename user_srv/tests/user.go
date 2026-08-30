package main

import (
	"context"
	"fmt"
	"mxshop_srvs/user_srv/proto"

	"google.golang.org/grpc"
)

var userClient proto.UserClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:8080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	userClient = proto.NewUserClient(conn)
}

// 测试获取数据并分页
func TestGetUserList() {
	rsp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    1,
		PSize: 9,
	})
	if err != nil {
		panic(err)
	}

	for _, user := range rsp.Data {
		fmt.Println(user.NickName, user.Mobile, user.PassWord)
		checkRespone, err := userClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
			Password:          "admin123",
			EncryptedPassword: user.PassWord,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(checkRespone.Success)
	}
}

func TestCreateUser() {
	for i := 0; i < 10; i++ {
		rsp, err := userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
			NickName: fmt.Sprintf("boddy%d", i),
			Mobile:   fmt.Sprintf("1878222222%d", i),
			Password: fmt.Sprintf("admin123"),
		})

		if err != nil {
			panic(err)
		}
		fmt.Println(rsp.Id)
	}
}

func TestUpdateUser() {
	rsp, err := userClient.UpdateUser(context.Background(), &proto.UpdateUserInfo{
		Id:       1,
		NickName: "jjj",
		Gender:   "man",
		BirtDay:  0,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(rsp)

}
func main() {
	Init()
	defer conn.Close()
	TestGetUserList()
	//TestCreateUser()
	TestUpdateUser()
}
