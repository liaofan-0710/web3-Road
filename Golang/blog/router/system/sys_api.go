package system

import (
	v1 "Project/api/v1"
	"github.com/gin-gonic/gin"
)

type ApiRouter struct{}

func (s *ApiRouter) InitApiRouter(Router *gin.RouterGroup) {
	// apiRouter := Router.Group("api").Use(middleware.Register())
	apiRouterWithoutRecord := Router.Group("api")
	apiRouterApi := v1.ApiGroupApp.SystemApiGroup.SystemApiApi
	{
		apiRouterWithoutRecord.POST("login", apiRouterApi.Login)   // 登录接口
		apiRouterWithoutRecord.POST("enroll", apiRouterApi.Enroll) // 注册接口
	}
}
