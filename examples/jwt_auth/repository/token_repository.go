package repository

import (
	"database/sql"
	"time"

	"project/examples/jwt_auth/model"
	"project/pkg/logger"
	"project/pkg/utils/idgen"
)

// SaveRefreshToken 保存 refresh token
func SaveRefreshToken(db *sql.DB, token *model.RefreshToken) error {

	token.ID, _ = idgen.GenerateUUID(16)
	query := `
		INSERT INTO refresh_tokens (id, user_id, token, expires_at, created_at, is_revoked)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := db.Exec(query,
		token.ID,
		token.UserID,
		token.Token,
		token.ExpiresAt,
		time.Now(),
		false,
	)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	return nil
}

// GetRefreshToken 获取 refresh token
func GetRefreshToken(db *sql.DB, token string) (*model.RefreshToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at, is_revoked
		FROM refresh_tokens WHERE token = ?
	`
	var rt model.RefreshToken
	err := db.QueryRow(query, token).Scan(
		&rt.ID,
		&rt.UserID,
		&rt.Token,
		&rt.ExpiresAt,
		&rt.CreatedAt,
		&rt.IsRevoked,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}
	return &rt, nil
}

// RevokeRefreshToken 撤销 refresh token
func RevokeRefreshToken(db *sql.DB, token string) error {
	query := `UPDATE refresh_tokens SET is_revoked = true WHERE token = ?`
	_, err := db.Exec(query, token)
	return err
}

// RevokeAllUserTokens 撤销用户所有 token（用于登出所有设备）
func RevokeAllUserTokens(db *sql.DB, userID string) error {
	query := `UPDATE refresh_tokens SET is_revoked = true WHERE user_id = ?`
	_, err := db.Exec(query, userID)
	return err
}

// DeleteExpiredTokens 删除过期的 token（定期清理）
func DeleteExpiredTokens(db *sql.DB) error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < ?`
	_, err := db.Exec(query, time.Now())
	return err
}
