package system

import "Project/service"

type ApiGroup struct {
	SystemApiApi
	JwtApi
	PostApi
}

var (
	apiService  = service.ServiceGroupApp.SystemServiceGroup.ApiService
	jwtService  = service.ServiceGroupApp.SystemServiceGroup.JwtService
	postService = service.ServiceGroupApp.SystemServiceGroup.PostService
)
