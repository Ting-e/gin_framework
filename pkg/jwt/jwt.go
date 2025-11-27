package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 错误定义
var (
	ErrTokenExpired     = errors.New("token expired")
	ErrTokenInvalid     = errors.New("token invalid")
	ErrTokenMalformed   = errors.New("token malformed")
	ErrTokenNotValidYet = errors.New("token not valid yet")
)

// Claims JWT 载荷
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// JWTConfig JWT 配置
type JWT struct {
	Secret           []byte        // 密钥
	Issuer           string        // 签发者
	ExpiresAt        time.Duration // 过期时间
	RefreshExpiresAt time.Duration // 刷新令牌过期时间
}

var defaultConfig = &JWT{
	Secret:           []byte("dev-secret-please-change-in-production"),
	Issuer:           "my-app",
	ExpiresAt:        2 * time.Hour,      // 访问令牌 2 小时
	RefreshExpiresAt: 7 * 24 * time.Hour, // 刷新令牌 7 天
}

// SetSecret 设置密钥
func SetSecret(secret string) {
	defaultConfig.Secret = []byte(secret)
}

// SetConfig 设置完整配置
func SetConfig(cfg *JWT) {
	if cfg.Secret != nil {
		defaultConfig.Secret = cfg.Secret
	}
	if cfg.Issuer != "" {
		defaultConfig.Issuer = cfg.Issuer
	}
	if cfg.ExpiresAt > 0 {
		defaultConfig.ExpiresAt = cfg.ExpiresAt
	}
	if cfg.RefreshExpiresAt > 0 {
		defaultConfig.RefreshExpiresAt = cfg.RefreshExpiresAt
	}
}

// GenerateToken 生成访问令牌
func GenerateToken(userID, username, role string) (string, error) {
	return GenerateTokenWithExpire(userID, username, role, defaultConfig.ExpiresAt)
}

// GenerateTokenWithExpire 生成指定过期时间的令牌
func GenerateTokenWithExpire(userID, username, role string, expiresAt time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    defaultConfig.Issuer,
			ExpiresAt: jwt.NewNumericDate(now.Add(expiresAt)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(defaultConfig.Secret)
}

// GenerateRefreshToken 生成刷新令牌
func GenerateRefreshToken(userID, username, role string) (string, error) {
	return GenerateTokenWithExpire(userID, username, role, defaultConfig.RefreshExpiresAt)
}

// GenerateTokenPair 生成访问令牌和刷新令牌
func GenerateTokenPair(userID, username, role string) (accessToken, refreshToken string, err error) {
	accessToken, err = GenerateToken(userID, username, role)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = GenerateRefreshToken(userID, username, role)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ParseToken 解析令牌
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenInvalid
		}
		return defaultConfig.Secret, nil
	})

	if err != nil {
		// 详细的错误处理
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrTokenMalformed
		}
		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, ErrTokenNotValidYet
		}
		return nil, ErrTokenInvalid
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// RefreshToken 刷新令牌（使用刷新令牌生成新的访问令牌）
func RefreshToken(refreshTokenString string) (string, error) {
	claims, err := ParseToken(refreshTokenString)
	if err != nil {
		return "", err
	}

	// 生成新的访问令牌
	return GenerateToken(claims.UserID, claims.Username, claims.Role)
}

// ValidateToken 验证令牌是否有效
func ValidateToken(tokenString string) bool {
	_, err := ParseToken(tokenString)
	return err == nil
}
