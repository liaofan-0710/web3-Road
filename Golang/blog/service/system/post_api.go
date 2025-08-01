package system

import (
	"Project/global"
	"Project/model"
	"Project/model/request"
)

type PostService struct{}

var ApiPostService = new(PostService)

// CreatePost @function: CreatePost
// @description: 新增基础api
// @param: api model.User
// @return: err error
func (postService PostService) CreatePost(post model.Post) (err error) {
	return global.BG_DB.Model(model.Post{}).Create(&post).Error
}

// QueryPost @function: QueryPost
// @description: 新增基础api
// @param: api model.Post
// @return: err error
func (postService PostService) QueryPost() (post []model.Post, err error) {
	err = global.BG_DB.Model(model.Post{}).Find(&post).Error
	return post, err
}

// GetPost @function: GetPost
// @description: 新增基础api
// @param: api model.Post
// @return: err error
func (postService PostService) GetPost(id uint) (post model.Post, err error) {
	err = global.BG_DB.Model(model.Post{}).Where("id = ?", id).First(&post).Error
	return post, err
}

// UpdatePost @function: UpdatePost
// @description: 新增基础api
// @param: api model.Post
// @return: err error
func (postService PostService) UpdatePost(post request.PostUpdate, userId uint) (err error) {
	return global.BG_DB.Model(model.Post{}).
		Where("user_id = ?", userId).Where("id = ?", post.ID).
		Update("title", post.Title).Update("content", post.Content).Error
}

// DeletePost @function: DeletePost
// @description: 新增基础api
// @param: api model.Post
// @return: err error
func (postService PostService) DeletePost(postId, userId uint) (err error) {
	return global.BG_DB.Model(model.Post{}).Where("id = ?", postId).Where("user_id = ?", userId).
		Delete(&model.Post{}).Error
}
