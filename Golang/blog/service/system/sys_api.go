package system

import (
	"Project/global"
	"Project/model"
)

type ApiService struct{}

var ApiServiceApp = new(ApiService)

// GetUsers @function: GetUsers
// @description: 新增基础api
// @param: api model.User
// @return: err error
func (apiService ApiService) GetUsers(username string) (user model.User, err error) {
	err = global.BG_DB.Where("username = ?", username).First(&user).Error
	return user, err
}

// CreateUsers @function: CreateUsers
// @description: 新增基础api
// @param: api model.User
// @return: err error
func (apiService ApiService) CreateUsers(user model.User) (err error) {
	return global.BG_DB.Create(&user).Error
}
