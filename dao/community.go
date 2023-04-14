package dao

import (
	"go.uber.org/zap"
	"humble_blog/config"
	"humble_blog/model"
)

func GetCommunityCategory() []model.Category {
	db := config.GetDB()
	var categorys []model.Category
	result := db.Find(&categorys)
	if err := result.Error; err != nil {
		zap.L().Error(err.Error())
	}
	return categorys
}
func GetCategoryByID(categoryId int32) model.Category {
	db := config.GetDB()
	var category model.Category
	result := db.Where(categoryId).Find(&category)
	if err := result.Error; err != nil {
		zap.L().Error(err.Error())
	}
	return category
}
