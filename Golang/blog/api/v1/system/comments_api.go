package system

import (
	"Project/model"
	"Project/model/request"
	"Project/model/response"
	"Project/utils"
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
	baseClaims, err := utils.GetContextUserInfo(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
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
	if err = commentsService.CreateComment(createCom); err != nil {
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

	comRes, err := commentsService.GetComment(comReq.PostID)
	if err != nil {
		response.FailWithMessage("error: "+err.Error(), c)
		return
	}

	response.OkWithData(comRes, c)
}
