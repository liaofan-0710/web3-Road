package system

import (
	"Project/global"
	"Project/model"
	"Project/model/request"
	"Project/model/response"
	"github.com/gin-gonic/gin"
)

type CommentsApi struct{}

// CommentsCreate
// @Tags      CommentsApi
// @Summary   评论创建
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.PostApi                  true  "api路径, api中文描述, api组, 方法"
// @Success   200   {object}  response.Response{msg=string}  "创建基础api"
// @Router    /api/enroll [post]
func (p CommentsApi) CommentsCreate(c *gin.Context) {
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

	var com request.Comment
	if err := c.ShouldBindJSON(&com); err != nil {
		response.FailWithMessage("error: "+err.Error(), c)
		return
	}
	createCom := model.Comment{
		PostID:  com.PostID,
		Content: com.Content,
		UserID:  baseClaims.ID,
	}

	if err := global.BG_DB.Model(model.Comment{}).Create(&createCom).Error; err != nil {
		response.FailWithMessage("error: "+err.Error(), c)
		return
	}

	response.Ok(c)
}

// CommentsQuery
// @Tags      CommentsApi
// @Summary   文章创建
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.PostApi                  true  "api路径, api中文描述, api组, 方法"
// @Success   200   {object}  response.Response{msg=string}  "创建基础api"
// @Router    /api/enroll [post]
func (p CommentsApi) CommentsQuery(c *gin.Context) {
	var comReq request.CommentQuery
	if err := c.ShouldBindJSON(&comReq); err != nil {
		response.FailWithMessage("error: "+err.Error(), c)
		return
	}
	var comRes []response.Comment
	if err := global.BG_DB.Model(model.Comment{}).Where("post_id = ?", comReq.PostID).Find(&comRes).Error; err != nil {
		response.FailWithMessage("error: "+err.Error(), c)
		return
	}

	response.OkWithData(comRes, c)
}
