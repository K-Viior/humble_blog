package logic

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
	"humble_blog/common"
	"humble_blog/config"
	"humble_blog/dao"
	"humble_blog/dto"
	"humble_blog/model"
	"strconv"
	"time"
)

// 创建帖子
func CreatePost(userId int64, postDto dto.PostDTO) {
	rdb := config.RDB
	pipeline := rdb.Pipeline()
	ctx := context.Background()
	//封装post
	post := model.Post{
		AuthorID:   userId,
		CategoryID: postDto.CategoryID,
		Title:      postDto.Title,
		Content:    postDto.Content,
		Status:     postDto.Status,
		BaseModel:  model.BaseModel{CreateTime: time.Now(), UpdateTime: time.Now()},
	}
	dao.CreatePost(&post)
	//redis记录帖子的发布时间
	pipeline.ZAdd(ctx, common.KeyPostTimeZSet, &redis.Z{
		Score:  float64(post.CreateTime.Unix()),
		Member: post.ID,
	})
	//记录帖子的分数为发布时间
	pipeline.ZAdd(ctx, common.KeyPostScoreZSet, &redis.Z{
		Score:  float64(post.CreateTime.Unix()),
		Member: post.ID,
	})
	if cmds, err := pipeline.Exec(ctx); err != nil {
		zap.L().Error("create post redis err"+err.Error(), zap.Any("cmds", cmds))
		return
	}
	zap.L().Info("redis记录发帖时间 " + fmt.Sprintf("%d", post.CreateTime.Unix()))
}

// 分页查询
func GetPostList(PLDto *dto.PostListQuery) *[]model.Post {
	var posts []model.Post
	//检查分页参数合法性
	CheckPageParams(PLDto)
	//进行分页查询
	switch PLDto.Order {
	case "score":
		posts = GetHotPostList(PLDto)
		return &posts
	case "create_time":
		posts = dao.GetPostList(PLDto)
	}
	ids := make([]string, 0)
	for _, post := range posts {
		ids = append(ids, strconv.Itoa(int(post.ID)))
	}
	ups, downs := GetVotes(ids)
	for i, _ := range posts {
		posts[i].Up = ups[i]
		posts[i].Down = downs[i]
	}
	return &posts
}

// 从redis获取热点文章
func GetHotPostList(plDto *dto.PostListQuery) []model.Post {
	rdb := config.RDB
	start := (plDto.Page - 1) * plDto.PageSize
	end := start + plDto.PageSize - 1
	//从redis中根据分页参数获取前?个文章id
	ids, err := rdb.ZRevRange(context.Background(),
		common.KeyPostScoreZSet,
		int64(start),
		int64(end),
	).Result()
	if err != nil {
		zap.L().Error(err.Error())
	}
	//根据文章Id获取文章
	posts := GetPostListInIDs(ids)
	//获取redis中的投票数据
	ups, downs := GetVotes(ids)
	for i, _ := range posts {
		posts[i].Up = ups[i]
		posts[i].Down = downs[i]
	}

	return posts
}
func GetVotes(ids []string) ([]int, []int) {
	rdb := config.RDB
	ctx := context.Background()
	ups := make([]int, 0)
	downs := make([]int, 0)
	for _, postID := range ids {
		result, _ := rdb.HVals(ctx, common.KeyPostVotedPrefix+postID).Result()
		var likes, dislikes int
		for _, vote := range result {
			if vote == "1" {
				likes++
			} else if vote == "-1" {
				dislikes++
			}
		}
		ups = append(ups, likes)
		downs = append(downs, dislikes)
	}
	return ups, downs
}

// 根据文章Id获取文章
func GetPostListInIDs(ids []string) []model.Post {
	db := config.GetDB()
	posts := make([]model.Post, 0)
	db.Where("id in ?", ids).Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "FIELD(id,?)", Vars: []any{ids}, WithoutParentheses: true}}).
		Find(&posts)
	return posts
}

func CheckPageParams(PLDto *dto.PostListQuery) {
	if PLDto.Page == 0 {
		PLDto.Page = 1
	}
	switch {
	case PLDto.PageSize > 100:
		PLDto.PageSize = 100
	case PLDto.PageSize < 0:
		PLDto.PageSize = 10
	}
}

// 获取文章详情
func GetPostById(postId string) *model.Post {
	//格式转换
	post_id, _ := strconv.ParseInt(postId, 10, 32)

	post := dao.GetPostById(int32(post_id))

	return post
}

const (
	Score = 666
)

// 点踩的方法
func PostVote(voteDTO *dto.VoteDTO, userId string) {
	postID := strconv.FormatInt(int64(voteDTO.PostID), 10)
	rdb := config.RDB
	pipeline := rdb.Pipeline()
	ctx := context.Background()
	//判断当前用户是否已经投票
	val := rdb.HGet(ctx, common.KeyPostVotedPrefix+postID, userId)
	record, err := strconv.Atoi(val.Val())
	if err != nil {
		zap.L().Error(err.Error())
	}
	//计算差值
	diff := voteDTO.Type - record
	//更新当前帖子的分值
	pipeline.ZIncrBy(ctx, common.KeyPostScoreZSet, float64(diff*Score), postID)
	//未投票则记录用户的投票
	pipeline.HSet(ctx, common.KeyPostVotedPrefix+postID, userId, strconv.Itoa(voteDTO.Type))

	if cmds, err := pipeline.Exec(ctx); err != nil {
		zap.L().Error("postVoting in redis err "+err.Error(), zap.Any("cmds", cmds))
	}
}
