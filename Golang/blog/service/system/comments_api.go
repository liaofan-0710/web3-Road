package system

import (
	"Project/global"
	"Project/model"
	"Project/model/response"
)

type CommentsService struct{}

var ApiCommentsService = new(CommentsService)

// CreateComment @function: CreateComment
// @description: 新增基础api
// @param: api model.Comment
// @return: err error
func (commentService CommentsService) CreateComment(com model.Comment) (err error) {
	return global.BG_DB.Model(model.Comment{}).Create(&com).Error
}

// GetComment @function: GetComment
// @description: 新增基础api
// @param: api model.Comment
// @return: err error
func (commentService CommentsService) GetComment(postId uint) (comRes []response.Comment, err error) {
	err = global.BG_DB.Model(model.Comment{}).Where("post_id = ?", postId).Find(&comRes).Error
	return comRes, err
}
