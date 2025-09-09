package system

import (
	v1 "Project/api/v1"
	"Project/middleware"
	"github.com/gin-gonic/gin"
)

type CommentsRouter struct{}

func (s *CommentsRouter) OperateCommentsRouter(Router *gin.RouterGroup) {
	postRouter := Router.Group("api").Use(middleware.JWTAuth())
	postRouterWithoutRecord := Router.Group("api")
	apiRouterApi := v1.ApiGroupApp.SystemApiGroup.CommentsApi
	{
		postRouter.POST("create-comments", apiRouterApi.CommentsCreate) // 评论创建功能
	}
	{
		postRouterWithoutRecord.GET("query-comments", apiRouterApi.CommentsQuery) // 单篇评论查询功能
	}
}
