package handler

import (
	"net/http"
	"project/examples/jwt_auth/model"
	srv "project/examples/jwt_auth/service"
	"project/pkg/json"
	"project/pkg/logger"
	"project/pkg/response"

	"github.com/gin-gonic/gin"
)

var service srv.APIService

func init() {
	service = srv.GetService()
}

// Register 用户注册
// @Summary 用户注册
// @Description 创建新用户账号
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body model.RegisterRequest true "注册信息"
// @Success 200 {object} model.AuthResponse
// @Router /api/auth/register [post]
func Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Error(err)
		response.BadRequest(c)
		return
	}

	logger.Sugar.Info("Register 请求参数：", string(json.JSONMarshal(req)))

	authResp, err := service.Register(&req)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Sugar.Info("Register 响应参数：", string(json.JSONMarshal(authResp)))

	response.Success(c, authResp)
}

// Login 用户登录
// @Summary 用户登录
// @Description 使用用户名和密码登录
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "登录信息"
// @Success 200 {object} model.AuthResponse
// @Router /api/auth/login [post]
func Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c)
		return
	}

	logger.Sugar.Info("Login 请求参数：", string(json.JSONMarshal(req)))

	authResp, err := service.Login(&req)
	if err != nil {
		response.Failed(c, http.StatusUnauthorized, err.Error())
		return
	}

	logger.Sugar.Info("Login 响应参数：", string(json.JSONMarshal(authResp)))

	response.Success(c, authResp)
}

// RefreshToken 刷新访问令牌
// @Summary 刷新访问令牌
// @Description 使用 refresh token 获取新的 access token
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body model.RefreshTokenRequest true "Refresh Token"
// @Success 200 {object} model.TokenPair
// @Router /api/auth/refresh [post]
func RefreshToken(c *gin.Context) {
	var req model.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c)
		return
	}

	logger.Sugar.Info("Login 请求参数：", string(json.JSONMarshal(req)))

	tokenPair, err := service.RefreshToken(req.RefreshToken)
	if err != nil {
		response.Failed(c, http.StatusUnauthorized, err.Error())
		return
	}

	logger.Sugar.Info("Login 响应参数：", string(json.JSONMarshal(tokenPair)))

	response.Success(c, tokenPair)
}

// Logout 登出
// @Summary 用户登出
// @Description 撤销 refresh token
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body model.RefreshTokenRequest true "Refresh Token"
// @Success 200 {string} string "登出成功"
// @Router /api/auth/logout [post]
func Logout(c *gin.Context) {
	var req model.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c)
		return
	}

	logger.Sugar.Info("Logout 请求参数：", string(json.JSONMarshal(req)))

	if err := service.Logout(req.RefreshToken); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	logger.Sugar.Info("Logout 响应参数：", string(json.JSONMarshal("登出成功")))

	response.Success(c, "登出成功")
}

// LogoutAllDevices 登出所有设备
// @Summary 登出所有设备
// @Description 撤销用户所有设备的 refresh token
// @Tags 认证
// @Security BearerAuth
// @Produce json
// @Success 200 {string} string "已登出所有设备"
// @Router /api/auth/logout-all [post]
func LogoutAllDevices(c *gin.Context) {
	// 从 JWT 中间件获取用户 ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}

	logger.Sugar.Info("LogoutAllDevices 请求参数：", string(json.JSONMarshal(userID)))

	if err := service.LogoutAllDevices(userID.(string)); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	logger.Sugar.Info("LogoutAllDevices 响应参数：", string(json.JSONMarshal("已登出所有设备")))

	response.Success(c, "已登出所有设备")
}

// GetUserInfo 获取当前用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户的信息
// @Tags 认证
// @Security BearerAuth
// @Produce json
// @Success 200 {object} model.UserInfo
// @Router /api/auth/userinfo [get]
func GetUserInfo(c *gin.Context) {
	// 从 JWT 中间件获取用户 ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}

	logger.Sugar.Info("GetUserInfo 请求参数：", string(json.JSONMarshal("userID")))

	userInfo, err := service.GetUserInfo(userID.(string))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	logger.Sugar.Info("GetUserInfo 响应参数：", string(json.JSONMarshal(userInfo)))

	response.Success(c, userInfo)
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改当前用户密码
// @Tags 认证
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.ChangePasswordRequest true "密码信息"
// @Success 200 {string} string "密码修改成功"
// @Router /api/auth/change-password [post]
func ChangePassword(c *gin.Context) {
	// 从 JWT 中间件获取用户 ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}

	var req model.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c)
		return
	}

	req.UserID = userID.(string)

	logger.Sugar.Info("ChangePassword 请求参数：", string(json.JSONMarshal(req)))

	if err := service.ChangePassword(&req); err != nil {
		response.Failed(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Sugar.Info("ChangePassword 响应参数：", string(json.JSONMarshal("密码修改成功，请重新登录")))

	response.Success(c, "密码修改成功，请重新登录")
}
