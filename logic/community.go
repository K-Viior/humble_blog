package logic

import (
	"github.com/gin-gonic/gin"
	"humble_blog/common"
	"humble_blog/dao"
)

// 获取社区分类的方法
func GetCommunityCategory(ctx *gin.Context) {
	categorys := dao.GetCommunityCategory()
	common.Success(ctx, categorys)
}
