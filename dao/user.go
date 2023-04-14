package dao

import (
	"errors"
	"go.uber.org/zap"
	"humble_blog/config"
	"humble_blog/model"
)

// 判断用户是否存在的方法
func IsUserExist(username string) bool {
	db := config.GetDB()
	var user model.User
	result := db.Where(&model.User{Username: username}).Find(&user)
	if result.RowsAffected == 0 {
		return false
	}
	return true
}

// 新增用户的方法
func Register(user *model.User) {
	db := config.GetDB()
	result := db.Create(&user)
	if err := result.Error; err != nil {
		panic(errors.New("创建新用户失败"))
	}
}

// 根据用户名查找用户
func GetUserByName(username string) *model.User {
	db := config.GetDB()
	var user model.User
	result := db.Where(&model.User{Username: username}).Find(&user)
	if err := result.Error; result.RowsAffected == 0 {
		zap.L().Error("Find User By UserName Error" + err.Error())
	}
	return &user
}
func GetUserById(userId int64) model.User {
	db := config.GetDB()
	var user model.User
	result := db.Where(&model.User{UserID: userId}).Find(&user)
	if err := result.Error; err != nil {
		zap.L().Error(err.Error())
	}
	return user
}
