// examples/jwt_auth/service/auth_service.go
package service

import (
	"errors"
	"time"

	"project/examples/jwt_auth/model"
	"project/examples/jwt_auth/repository"
	"project/pkg/config"
	"project/pkg/database"
	"project/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

// Register 用户注册
func (s *Service) Register(req *model.RegisterRequest) (*model.AuthResponse, error) {

	db := database.GetMysql().GetDb()

	// 1. 检查用户名是否已存在
	existingUser, err := repository.GetByUsername(db, req.Username)
	if err != nil && err.Error() != "用户不存在" {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("用户名已存在")
	}

	// 2. 检查邮箱是否已存在
	existingEmail, err := repository.GetByEmail(db, req.Email)
	if err != nil {
		return nil, err
	}
	if existingEmail != nil {
		return nil, errors.New("邮箱已被注册")
	}

	// 3. 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 4. 创建用户
	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         "user", // 默认角色
	}

	if err := repository.CreateUser(db, user); err != nil {
		return nil, err
	}

	// 5. 生成 token 并返回
	return s.generateAuthResponse(user)
}

// Login 用户登录
func (s *Service) Login(req *model.LoginRequest) (*model.AuthResponse, error) {

	db := database.GetMysql().GetDb()

	// 1. 查找用户
	user, err := repository.GetByUsername(db, req.Username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 2. 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 3. 生成 token 并返回
	return s.generateAuthResponse(user)
}

// RefreshToken 刷新 access token
func (s *Service) RefreshToken(refreshTokenStr string) (*model.TokenPair, error) {

	db := database.GetMysql().GetDb()

	jwtConfig := config.Get().JWT

	// 1. 验证 refresh token 格式
	_, err := jwt.ParseToken(refreshTokenStr)
	if err != nil {
		return nil, errors.New("无效的 refresh token")
	}

	// 2. 从数据库验证 refresh token
	storedToken, err := repository.GetRefreshToken(db, refreshTokenStr)
	if err != nil {
		return nil, err
	}
	if storedToken == nil {
		return nil, errors.New("refresh token 不存在")
	}

	// 3. 检查是否已撤销
	if storedToken.IsRevoked {
		return nil, errors.New("refresh token 已失效")
	}

	// 4. 检查是否过期
	if time.Now().After(storedToken.ExpiresAt) {
		return nil, errors.New("refresh token 已过期")
	}

	// 5. 获取用户信息
	user, err := repository.GetByID(db, storedToken.UserID)
	if err != nil {
		return nil, err
	}

	// 6. 生成新的 access token
	accessToken, err := jwt.GenerateToken(
		user.ID,
		user.Username,
		user.Role,
	)
	if err != nil {
		return nil, err
	}

	// 7. 可选：轮换 refresh token（更安全）
	// 撤销旧的 refresh token
	if err := repository.RevokeRefreshToken(db, refreshTokenStr); err != nil {
		return nil, err
	}

	// 生成新的 refresh token
	newRefreshToken, err := jwt.GenerateRefreshToken(
		user.ID,
		user.Username,
		user.Role,
	)
	if err != nil {
		return nil, err
	}

	// 保存新的 refresh token
	rt := &model.RefreshToken{
		UserID:    user.ID,
		Token:     newRefreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * time.Duration(jwtConfig.RefreshExpiresDays)),
	}
	if err := repository.SaveRefreshToken(db, rt); err != nil {
		return nil, err
	}

	return &model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// Logout 登出（撤销 refresh token）
func (s *Service) Logout(refreshToken string) error {
	db := database.GetMysql().GetDb()
	return repository.RevokeRefreshToken(db, refreshToken)
}

// LogoutAllDevices 登出所有设备
func (s *Service) LogoutAllDevices(userID string) error {
	db := database.GetMysql().GetDb()
	return repository.RevokeAllUserTokens(db, userID)
}

// GetUserInfo 获取用户信息
func (s *Service) GetUserInfo(userID string) (*model.UserInfo, error) {
	db := database.GetMysql().GetDb()
	user, err := repository.GetByID(db, userID)
	if err != nil {
		return nil, err
	}

	return &model.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}, nil
}

// ChangePassword 修改密码
func (s *Service) ChangePassword(req *model.ChangePasswordRequest) error {
	db := database.GetMysql().GetDb()
	// 1. 获取用户
	user, err := repository.GetByID(db, req.UserID)
	if err != nil {
		return err
	}

	// 2. 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		return errors.New("原密码错误")
	}

	// 3. 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 4. 更新密码
	if err := repository.UpdatePassword(db, req.UserID, string(hashedPassword)); err != nil {
		return err
	}

	// 5. 撤销所有 refresh token（强制重新登录）
	return repository.RevokeAllUserTokens(db, req.UserID)
}

// generateAuthResponse 生成认证响应（内部方法）
func (s *Service) generateAuthResponse(user *model.User) (*model.AuthResponse, error) {

	db := database.GetMysql().GetDb()

	jwtConfig := config.Get().JWT

	// 1. 生成 access token
	accessToken, err := jwt.GenerateToken(
		user.ID,
		user.Username,
		user.Role,
	)
	if err != nil {
		return nil, err
	}

	// 2. 生成 refresh token
	refreshToken, err := jwt.GenerateRefreshToken(
		user.ID,
		user.Username,
		user.Role,
	)
	if err != nil {
		return nil, err
	}

	// 3. 保存 refresh token 到数据库
	rt := &model.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * time.Duration(jwtConfig.RefreshExpiresDays)),
	}
	if err := repository.SaveRefreshToken(db, rt); err != nil {
		return nil, err
	}

	// 4. 返回响应
	return &model.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(jwtConfig.ExpiresHours * 3600), // 转换为秒
		User: &model.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		},
	}, nil
}
