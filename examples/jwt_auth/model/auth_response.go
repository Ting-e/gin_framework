package model

// AuthResponse 认证响应
type AuthResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"` // Bearer
	ExpiresIn    int64     `json:"expires_in"` // 秒
	User         *UserInfo `json:"user"`
}

// UserInfo 用户信息（不包含敏感信息）
type UserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

// TokenPair Token 对
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
