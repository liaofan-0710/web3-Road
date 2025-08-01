package system

import (
	"Project/global"
	"Project/model"
	"Project/model/request"
	"Project/model/response"
	"Project/utils"
	"github.com/gin-gonic/gin"
)

type PostApi struct{}

// PostCreate
// @Tags      PostApi
// @Summary   文章创建
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.PostApi                  true  "api路径, api中文描述, api组, 方法"
// @Success   200   {object}  response.Response{msg=string}  "创建基础api"
// @Router    /api/enroll [post]
func (p PostApi) PostCreate(c *gin.Context) {
	baseClaims, err := utils.GetContextUserInfo(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	var post request.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		response.FailWithMessage("error: "+err.Error(), c)
		return
	}
	createPost := model.Post{
		Title:   post.Title,
		Content: post.Content,
		UserID:  baseClaims.ID,
	}

	if err := postService.CreatePost(createPost); err != nil {
		response.FailWithMessage("error: "+err.Error(), c)
		return
	}
	response.Ok(c)
}

// PostQuery
// @Tags      PostApi
// @Summary   文章创建
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.PostApi                  true  "api路径, api中文描述, api组, 方法"
// @Success   200   {object}  response.Response{msg=string}  "创建基础api"
// @Router    /api/enroll [post]
func (p PostApi) PostQuery(c *gin.Context) {
	postRes, err := postService.QueryPost()
	if err != nil {
		response.FailWithMessage("error: "+err.Error(), c)
		return
	}

	response.OkWithData(postRes, c)
}

// PostUpdate
// @Tags      PostApi
// @Summary   文章创建
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.PostApi                  true  "api路径, api中文描述, api组, 方法"
// @Success   200   {object}  response.Response{msg=string}  "创建基础api"
// @Router    /api/enroll [post]
func (p PostApi) PostUpdate(c *gin.Context) {
	baseClaims, err := utils.GetContextUserInfo(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	var post request.PostUpdate
	if err := c.ShouldBindJSON(&post); err != nil {
		response.FailWithMessage("error: "+err.Error(), c)
		return
	}
	// 判断该文章是否属于当前作者，否则返回无权限操作
	postModel, err := postService.GetPost(post.ID)
	if err != nil {
		response.FailWithMessage("查无此文章"+err.Error(), c)
		return
	}
	if err := global.BG_DB.Model(model.Post{}).Where("id = ?", post.ID).First(&postModel).Error; err != nil {

	}
	if postModel.UserID != baseClaims.ID {
		response.FailWithMessage("你无权限操作该文章", c)
		return
	}

	if err := postService.UpdatePost(post, baseClaims.ID); err != nil {
		response.FailWithMessage("error: "+err.Error(), c)
		return
	}

	response.Ok(c)
}

// PostDelete
// @Tags      PostApi
// @Summary   文章创建
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.PostApi                  true  "api路径, api中文描述, api组, 方法"
// @Success   200   {object}  response.Response{msg=string}  "创建基础api"
// @Router    /api/enroll [post]
func (p PostApi) PostDelete(c *gin.Context) {
	baseClaims, err := utils.GetContextUserInfo(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	var post request.PostDelete
	if err := c.ShouldBindJSON(&post); err != nil {
		response.FailWithMessage("error: "+err.Error(), c)
		return
	}
	postModel, err := postService.GetPost(post.ID)
	if err != nil {
		response.FailWithMessage("查无此文章"+err.Error(), c)
		return
	}
	if postModel.UserID != baseClaims.ID {
		response.FailWithMessage("你无权限操作该文章", c)
		return
	}
	if err = postService.DeletePost(post.ID, baseClaims.ID); err != nil {
		response.FailWithMessage("error: "+err.Error(), c)
		return
	}

	response.Ok(c)
}
