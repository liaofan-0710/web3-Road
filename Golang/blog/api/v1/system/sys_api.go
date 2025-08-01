package system

import (
	"Project/global"
	"Project/model"
	"Project/model/request"
	"Project/model/response"
	"Project/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type SysApiApi struct{}

// Enroll
// @Tags      SysApi
// @Summary   用户注册
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysApiApi                  true  "api路径, api中文描述, api组, 方法"
// @Success   200   {object}  response.Response{msg=string}  "创建基础api"
// @Router    /api/enroll [post]
func (s SysApiApi) Enroll(c *gin.Context) {
	var user request.EnrollUser
	if err := c.ShouldBindJSON(&user); err != nil {
		response.FailWithMessage("error: "+err.Error(), c)
		return
	}
	// 加密密码
	hashedPassword, err1 := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err1 != nil {
		response.FailWithMessage("error: Failed to hash password"+err1.Error(), c)
		return
	}
	userInfo := model.User{
		Username: user.Username,
		Password: string(hashedPassword),
		Email:    user.Email,
	}
	if err := apiService.CreateUsers(userInfo).Error; err != nil {
		response.FailWithMessage("error: Failed to create user", c)
		return
	}

	response.OkWithMessage("User registered successfully", c)
}

// Login
// @Tags      SysApiApi
// @Summary   用户登录接口
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysApiApi                  true  "api路径, api中文描述, api组, 方法"
// @Success   200   {object}  response.Response{msg=string}  "创建基础api"
// @Router    /api/login [post]
func (s SysApiApi) Login(c *gin.Context) {
	var user request.LoginUser
	if err := c.ShouldBindJSON(&user); err != nil {
		response.FailWithMessage("error"+err.Error(), c)
		return
	}

	storedUser, err := apiService.GetUsers(user.Username)
	if err != nil {
		response.FailWithMessage("error: "+err.Error(), c)
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
		response.FailWithMessage("error: Invalid username or password", c)
		return
	}

	//// 生成 JWT
	//token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	//	"id":       storedUser.ID,
	//	"username": storedUser.Username,
	//	"exp":      time.Now().Add(time.Hour * 24).Unix(),
	//})
	//
	//tokenString, err := token.SignedString([]byte("your_secret_key"))
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
	//	return
	//}
	s.TokenNext(c, storedUser)
}

// TokenNext 登录以后签发jwt
func (s SysApiApi) TokenNext(c *gin.Context, user model.User) {
	j := &utils.JWT{SigningKey: []byte(global.BG_CONFIG.JWT.SigningKey)} // 唯一签名
	claims := j.CreateClaims(request.BaseClaims{
		ID:       user.ID,
		Username: user.Username,
	})
	token, err := j.CreateToken(claims)
	if err != nil {
		global.BG_LOG.Error("获取token失败!", zap.Error(err))
		response.FailWithMessage("获取token失败", c)
		return
	}
	if !global.BG_CONFIG.System.UseMultipoint {
		response.OkWithDetailed(response.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
		}, "登录成功", c)
		return
	}

	if jwtStr, err := jwtService.GetRedisJWT(user.Username); err == redis.Nil {
		if err := jwtService.SetRedisJWT(token, user.Username); err != nil {
			global.BG_LOG.Error("设置登录状态失败!", zap.Error(err))
			response.FailWithMessage("设置登录状态失败", c)
			return
		}
		response.OkWithDetailed(response.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
		}, "登录成功", c)
	} else if err != nil {
		global.BG_LOG.Error("设置登录状态失败!", zap.Error(err))
		response.FailWithMessage("设置登录状态失败", c)
	} else {
		var blackJWT model.JwtBlacklist
		blackJWT.Jwt = jwtStr
		if err := jwtService.JsonInBlacklist(blackJWT); err != nil {
			response.FailWithMessage("jwt作废失败", c)
			return
		}
		if err := jwtService.SetRedisJWT(token, user.Username); err != nil {
			response.FailWithMessage("设置登录状态失败", c)
			return
		}
		response.OkWithDetailed(response.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
		}, "登录成功", c)
	}
}
