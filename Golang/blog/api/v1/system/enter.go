package system

import "Project/service"

type ApiGroup struct {
	SysApiApi
	JwtApi
	PostApi
	CommentsApi
}

var (
	apiService      = service.ServiceGroupApp.SystemServiceGroup.ApiService
	jwtService      = service.ServiceGroupApp.SystemServiceGroup.JwtService
	postService     = service.ServiceGroupApp.SystemServiceGroup.PostService
	commentsService = service.ServiceGroupApp.SystemServiceGroup.CommentsService
)
