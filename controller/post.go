package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"humble_blog/common"
	"humble_blog/dto"
	"humble_blog/logic"
	"net/http"
	"strconv"
)

// 创建帖子
func CreatePost(ctx *gin.Context) {
	var postDto dto.PostDTO
	//将前端传来的数据绑定
	err := ctx.ShouldBind(&postDto)
	if err != nil {
		zap.L().Error(err.Error())
		return
	}
	//获取用户ID,判断是否登录
	userId, exists := ctx.Get("userId")
	if !exists {
		zap.L().Info("用户未登录")
	}
	logic.CreatePost(userId.(int64), postDto)
	common.Success(ctx, nil)
}

// 分页查询帖子
func GetPostList(ctx *gin.Context) {
	postListQuery := dto.PostListQuery{
		Page:     1,
		PageSize: 10,
		Order:    "create_time",
	}
	if err := ctx.ShouldBind(&postListQuery); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	posts := logic.GetPostList(&postListQuery)
	common.Success(ctx, posts)
}

// 根据文章id查找文章内容
func GetPostById(ctx *gin.Context) {
	//从前端获取文章ID
	postId := ctx.Param("postId")
	if postId == "" {
		common.FailByMsg(ctx, "postId为空")
	}
	//获取文章内容
	post := logic.GetPostById(postId)

	common.Success(ctx, post)
}

func PostVote(ctx *gin.Context) {
	//从前端获取投票数据
	var voteDto dto.VoteDTO
	if err := ctx.ShouldBind(&voteDto); err != nil {
		zap.L().Error(err.Error())
	}
	userId, _ := ctx.Get("userId")
	//保存投票数据
	logic.PostVote(&voteDto, strconv.FormatInt(userId.(int64), 10))

	common.Success(ctx, nil)
}
