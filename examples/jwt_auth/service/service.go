package service

import (
	"project/examples/jwt_auth/model"
	"project/pkg/utils/snowflake"
)

var defaultService *Service

type Service struct {
	snowflake *snowflake.Worker
}

func GetService() *Service {
	return defaultService
}

func init() {
	defaultService = &Service{
		snowflake: snowflake.NewWorker(snowflake.WorkerID, snowflake.WataCenterID),
	}
}

type APIService interface {
	Register(req *model.RegisterRequest) (*model.AuthResponse, error)
	Login(req *model.LoginRequest) (*model.AuthResponse, error)
	RefreshToken(refreshTokenStr string) (*model.TokenPair, error)
	Logout(refreshToken string) error
	LogoutAllDevices(userID string) error
	GetUserInfo(userID string) (*model.UserInfo, error)
	ChangePassword(req *model.ChangePasswordRequest) error
	generateAuthResponse(user *model.User) (*model.AuthResponse, error)
}
