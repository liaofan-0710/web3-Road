package system

import (
	"Project/global"
	"Project/model"
	"Project/model/request"
	"Project/model/response"
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
	claims, ok := c.Get("claims")
	if !ok {
		response.FailWithMessage("user get fial", c)
		return
	}
	baseClaims, ok := claims.(*request.CustomClaims)
	if !ok {
		response.FailWithMessage("claims type mismatch", c)
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

	if err := global.BG_DB.Model(model.Post{}).Create(&createPost).Error; err != nil {
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
	var postRes []response.Post

	if err := global.BG_DB.Model(model.Post{}).Find(&postRes).Error; err != nil {
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
	claims, ok := c.Get("claims")
	if !ok {
		response.FailWithMessage("user get fial", c)
		return
	}
	baseClaims, ok := claims.(*request.CustomClaims)
	if !ok {
		response.FailWithMessage("claims type mismatch", c)
		return
	}

	var post request.PostUpdate
	var postModel model.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		response.FailWithMessage("error: "+err.Error(), c)
		return
	}
	// 判断该文章是否属于当前作者，否则返回无权限操作
	if err := global.BG_DB.Model(model.Post{}).Where("id = ?", post.ID).First(&postModel).Error; err != nil {
		response.FailWithMessage("查无此文章"+err.Error(), c)
		return
	}
	if postModel.UserID != baseClaims.ID {
		response.FailWithMessage("你无权限操作该文章", c)
		return
	}
	if err := global.BG_DB.Model(model.Post{}).
		Where("user_id = ?", baseClaims.ID).Where("id = ?", post.ID).
		Update("title", post.Title).Update("content", post.Content).Error; err != nil {
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
	claims, ok := c.Get("claims")
	if !ok {
		response.FailWithMessage("user get fial", c)
		return
	}
	baseClaims, ok := claims.(*request.CustomClaims)
	if !ok {
		response.FailWithMessage("claims type mismatch", c)
		return
	}

	var post request.PostDelete
	var postModel model.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		response.FailWithMessage("error: "+err.Error(), c)
		return
	}
	if err := global.BG_DB.Model(model.Post{}).Where("id = ?", post.ID).First(&postModel).Error; err != nil {
		response.FailWithMessage("查无此文章"+err.Error(), c)
		return
	}
	if postModel.UserID != baseClaims.ID {
		response.FailWithMessage("你无权限操作该文章", c)
		return
	}
	if err := global.BG_DB.Model(model.Post{}).Where("id = ?", post.ID).Where("user_id = ?", baseClaims.ID).
		Delete(&model.Post{}).Error; err != nil {
		response.FailWithMessage("error: "+err.Error(), c)
		return
	}

	response.Ok(c)
}
