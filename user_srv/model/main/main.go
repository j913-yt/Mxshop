package main

import (
	"crypto/md5"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	"github.com/anaskhan96/go-password-encoder"
)

func genMd5(code string) string {
	Md5 := md5.New()
	_, _ = io.WriteString(Md5, code)
	return hex.EncodeToString(Md5.Sum(nil))
}
func main() {
	//dsn := "root:123456@tcp(127.0.0.1:3306)/mxshop_user_srv?charset=utf8&parseTime=True&loc=Local"
	//newLogger := logger.New(
	//	log.New(os.Stdout, "\r\n", log.LstdFlags),
	//	logger.Config{
	//		SlowThreshold: time.Second,
	//		LogLevel:      logger.Info,
	//		Colorful:      true,
	//	},
	//)
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
	//	NamingStrategy: schema.NamingStrategy{
	//		SingularTable: true,
	//	},
	//	Logger: newLogger,
	//})

	//if err != nil {
	//	//	panic(err)
	//	//}
	//	//
	//	//_ = db.AutoMigrate(&model.User{})

	//fmt.Println(genMd5("123456"))

	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode("jytt password", options)
	newPassword := fmt.Sprintf("pbkdf2-sha512$%s$%s", salt, encodedPwd)
	fmt.Println(len(newPassword))
	fmt.Println(newPassword)

	passwordInfo := strings.Split(newPassword, "$")

	fmt.Println(passwordInfo)
	check := password.Verify("jytt password", passwordInfo[1], passwordInfo[2], options)
	fmt.Println(check)
}
