package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

const (
	testUserID   = "123"
	testUsername = "alice"
	testRole     = "admin"
)

func TestGenerateAndParseToken(t *testing.T) {
	// 重置配置避免污染
	defaultConfig.Secret = []byte("test-secret-for-unit-test")
	defaultConfig.ExpiresAt = 10 * time.Minute

	token, err := GenerateToken(testUserID, testUsername, testRole)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := ParseToken(token)
	assert.NoError(t, err)
	assert.Equal(t, testUserID, claims.UserID)
	assert.Equal(t, testUsername, claims.Username)
	assert.Equal(t, testRole, claims.Role)
}

func TestGenerateTokenPair(t *testing.T) {
	defaultConfig.Secret = []byte("pair-test-secret")
	defaultConfig.ExpiresAt = 5 * time.Minute
	defaultConfig.RefreshExpiresAt = 1 * time.Hour

	accessToken, refreshToken, err := GenerateTokenPair(testUserID, testUsername, testRole)
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)

	// 验证两个 token 内容一致（除过期时间）
	accessClaims, _ := ParseToken(accessToken)
	refreshClaims, _ := ParseToken(refreshToken)

	assert.Equal(t, accessClaims.UserID, refreshClaims.UserID)
	assert.Equal(t, accessClaims.Username, refreshClaims.Username)
	assert.Equal(t, accessClaims.Role, refreshClaims.Role)
}

func TestRefreshToken(t *testing.T) {
	defaultConfig.Secret = []byte("refresh-test-secret")
	defaultConfig.ExpiresAt = 1 * time.Minute
	defaultConfig.RefreshExpiresAt = 2 * time.Minute

	_, refreshToken, err := GenerateTokenPair(testUserID, testUsername, testRole)
	assert.NoError(t, err)

	newAccessToken, err := RefreshToken(refreshToken)
	assert.NoError(t, err)
	assert.NotEmpty(t, newAccessToken)

	// 验证新 token 可解析
	claims, err := ParseToken(newAccessToken)
	assert.NoError(t, err)
	assert.Equal(t, testUserID, claims.UserID)
}

func TestParseToken_Expired(t *testing.T) {
	// 手动构造一个已过期的 token
	expiredClaims := Claims{
		UserID:   testUserID,
		Username: testUsername,
		Role:     testRole,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    defaultConfig.Issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	tokenStr, _ := token.SignedString([]byte("test-secret-expired"))

	_, err := ParseToken(tokenStr)
	assert.ErrorIs(t, err, ErrTokenExpired)
}

func TestParseToken_Malformed(t *testing.T) {
	_, err := ParseToken("invalid.token.string")
	assert.ErrorIs(t, err, ErrTokenMalformed)
}

func TestParseToken_InvalidSignature(t *testing.T) {
	// 用正确结构但错误密钥签发
	correctToken, _ := GenerateToken(testUserID, testUsername, testRole)
	// 修改密钥后解析（模拟篡改）
	SetSecret("wrong-secret")

	_, err := ParseToken(correctToken)
	assert.ErrorIs(t, err, ErrTokenInvalid)

	// 恢复密钥
	SetSecret("dev-secret-please-change-in-production")
}

func TestParseToken_NotValidYet(t *testing.T) {
	futureClaims := Claims{
		UserID:   testUserID,
		Username: testUsername,
		Role:     testRole,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    defaultConfig.Issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)), // 10分钟后才生效
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, futureClaims)
	tokenStr, _ := token.SignedString(defaultConfig.Secret)

	_, err := ParseToken(tokenStr)
	assert.ErrorIs(t, err, ErrTokenNotValidYet)
}

func TestValidateToken(t *testing.T) {
	defaultConfig.Secret = []byte("validate-test-secret")

	token, _ := GenerateToken(testUserID, testUsername, testRole)
	assert.True(t, ValidateToken(token))
	assert.False(t, ValidateToken("invalid"))
}

func TestSetSecret(t *testing.T) {
	original := string(defaultConfig.Secret)
	SetSecret("new-test-secret")
	assert.Equal(t, "new-test-secret", string(defaultConfig.Secret))
	SetSecret(original) // restore
}

func TestSetConfig(t *testing.T) {
	originalSecret := defaultConfig.Secret
	originalIssuer := defaultConfig.Issuer
	originalExp := defaultConfig.ExpiresAt

	newCfg := &JWT{
		Secret:    []byte("config-test-secret"),
		Issuer:    "test-issuer",
		ExpiresAt: 30 * time.Second,
	}
	SetConfig(newCfg)

	assert.Equal(t, "config-test-secret", string(defaultConfig.Secret))
	assert.Equal(t, "test-issuer", defaultConfig.Issuer)
	assert.Equal(t, 30*time.Second, defaultConfig.ExpiresAt)

	// restore
	SetConfig(&JWT{
		Secret:    []byte(originalSecret),
		Issuer:    originalIssuer,
		ExpiresAt: originalExp,
	})
}
