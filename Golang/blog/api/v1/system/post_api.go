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
	var post request.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		response.FailWithMessage("error: "+err.Error(), c)
	}
	createPost := model.Post{
		Title:   post.Title,
		Content: post.Content,
	}

	if err := global.BG_DB.Debug().Model(model.Post{}).Create(&createPost).Error; err != nil {
		response.FailWithMessage("error: "+err.Error(), c)
	}

	response.Ok(c)
}

// QueryCreate
// @Tags      PostApi
// @Summary   文章创建
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.PostApi                  true  "api路径, api中文描述, api组, 方法"
// @Success   200   {object}  response.Response{msg=string}  "创建基础api"
// @Router    /api/enroll [post]
func (p PostApi) QueryCreate(c *gin.Context) {
	var post []model.Post
	var postRes []response.Post
	if err := global.BG_DB.Model(model.Post{}).Find(&post).Error; err != nil {
		response.FailWithMessage("error: "+err.Error(), c)
	}

	for _, v := range post {
		postRes = append(postRes, response.Post{Title: v.Title, Content: v.Content})
	}

	response.OkWithData(postRes, c)
}
