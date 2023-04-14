package dao

import (
	"go.uber.org/zap"
	"humble_blog/config"
	"humble_blog/dto"
	"humble_blog/model"
)

// 创建新贴的方法
func CreatePost(post *model.Post) {
	db := config.GetDB()
	result := db.Create(post)
	if result.RowsAffected == 0 {
		zap.L().Error("Create new Post failed")
	}
}

// 分页查询
func GetPostList(PLDto *dto.PostListQuery) []model.Post {
	db := config.GetDB()
	posts := make([]model.Post, 0)

	result := db.Preload("User").
		Preload("Category").
		Offset((PLDto.Page - 1) * PLDto.PageSize).
		Limit(PLDto.PageSize).
		Order(PLDto.Order + " DESC").
		Find(&posts)
	if err := result.Error; err != nil {
		zap.L().Error(err.Error())
	}
	return posts
}

// 获取文章详情
func GetPostById(post_id int32) *model.Post {
	db := config.GetDB()
	var post model.Post
	result := db.Where(post_id).Find(&post)
	if err := result.Error; err != nil {
		zap.L().Error(err.Error())
	}
	return &post
}
