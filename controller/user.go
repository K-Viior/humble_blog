package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"humble_blog/common"
	"humble_blog/dto"
	"humble_blog/logic"
	"net/http"
)

func Register(ctx *gin.Context) {
	var registerDto dto.RegisterDto
	//获得前端传来的用户参数
	if err := ctx.ShouldBind(&registerDto); err != nil {
		zap.L().Error(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	//保存用户信息
	logic.Register(&registerDto)
	common.Success(ctx, nil)
}
func Login(ctx *gin.Context) {
	//获取前端参数
	var userDto dto.UserDto
	if err := ctx.ShouldBindJSON(&userDto); err != nil {
		zap.L().Error(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	//执行用户登录，获取返回token
	tokens := logic.Login(&userDto, ctx)
	common.Success(ctx, tokens)
}
