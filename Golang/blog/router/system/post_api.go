package system

import (
	v1 "Project/api/v1"
	"Project/middleware"
	"github.com/gin-gonic/gin"
)

type PostRouter struct{}

func (s *PostRouter) OperatePostRouter(Router *gin.RouterGroup) {
	postRouter := Router.Group("api").Use(middleware.JWTAuth())
	postRouterWithoutRecord := Router.Group("api")
	apiRouterApi := v1.ApiGroupApp.SystemApiGroup.PostApi
	{
		postRouter.POST("create-post", apiRouterApi.PostCreate) // 文章创建功能
		postRouter.PUT("update-post", apiRouterApi.PostUpdate)  // 文章更新功能
		postRouter.DELETE("del-post", apiRouterApi.PostDelete)  // 文章删除功能
	}
	{
		postRouterWithoutRecord.GET("query-post", apiRouterApi.PostQuery) // 文章查询功能
	}
}
