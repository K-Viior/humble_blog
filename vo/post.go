package vo

import "humble_blog/model"

type PostDetail struct {
	AuthorName   string `json:"author_name"`
	CategoryName string `json:"category_name"`
	model.Post
}
