package logic

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"humble_blog/common"
	"humble_blog/config"
	"humble_blog/dao"
	"humble_blog/dto"
	"humble_blog/model"
	"humble_blog/pkg/jwt"
	"humble_blog/pkg/passwd"
	"humble_blog/pkg/snowflake"
	"humble_blog/vo"
	"strconv"
	"time"
)

func Register(registerDto *dto.RegisterDto) {
	//判断用户是否重复
	if dao.IsUserExist(registerDto.Username) {
		//用户存在，返回异常
		panic(common.NewCustomError(common.CodeUserExist))
	}
	//用户不存在，封装用户数据,注册用户
	dao.Register(&model.User{
		UserID:   snowflake.GenerateID(),
		Username: registerDto.Username,
		Password: passwd.Encode(registerDto.Password),
		Email:    registerDto.Email,
		Age:      registerDto.Age,
		BaseModel: model.BaseModel{
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
			DeleteTime: gorm.DeletedAt{},
		},
	})
}

// 用户登录的逻辑
func Login(userDto *dto.UserDto, ctx *gin.Context) vo.Tokens {

	//从数据库获取用户信息
	dbuser := dao.GetUserByName(userDto.Username)

	//验证user
	if dbuser.UserID == 0 {
		panic(common.NewCustomError(common.CodeInvalidPassword))
	}
	//验证密码
	if !passwd.Verify(userDto.Password, dbuser.Password) {
		panic(common.NewCustomError(common.CodeInvalidPassword))
	}
	//获取token
	accessToken, err := jwt.AccessToken(dbuser.UserID)
	if err != nil {
		zap.L().Error(err.Error())
	}
	refreshToken, err := jwt.RefreshToken(dbuser.UserID)
	if err != nil {
		zap.L().Error(err.Error())
	}
	//将token存入redis，做单设备登录校验（未完）
	// 将access_token 存入redis中 限制同一用户同一IP 同一时间只能登录一个设备
	// key user:token:user_id:IP value access_token
	config.RDB.Set(
		context.Background(),
		common.KeyUserTokenPrefix+strconv.FormatInt(dbuser.UserID, 10)+":"+ctx.RemoteIP(),
		accessToken, 2*time.Hour,
	)
	//返回token
	return vo.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
