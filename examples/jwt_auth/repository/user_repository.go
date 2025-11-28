package repository

import (
	"database/sql"
	"errors"
	"time"

	"project/examples/jwt_auth/model"
	"project/pkg/logger"
	"project/pkg/utils/idgen"
)

// CreateUser 创建用户
func CreateUser(db *sql.DB, user *model.User) error {

	user.ID, _ = idgen.GenerateUUID(16)

	query := `
		INSERT INTO user (user_id, username, email, password_hash, role, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.Exec(query,
		user.ID,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Role,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	return nil
}

// GetByUsername 根据用户名获取用户
func GetByUsername(db *sql.DB, username string) (*model.User, error) {
	query := `
		SELECT user_id, username, email, password_hash, role, created_at, updated_at
		FROM user WHERE username = ?
	`
	var user model.User
	err := db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("用户不存在")
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByID 根据 ID 获取用户
func GetByID(db *sql.DB, id string) (*model.User, error) {
	query := `
		SELECT user_id, username, email, password_hash, role, created_at, updated_at
		FROM user WHERE user_id = ?
	`
	var user model.User
	err := db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("用户不存在")
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func GetByEmail(db *sql.DB, email string) (*model.User, error) {
	query := `
		SELECT user_id, username, email, password_hash, role, created_at, updated_at
		FROM user WHERE email = ?
	`
	var user model.User
	err := db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdatePassword 更新密码
func UpdatePassword(db *sql.DB, userID string, newPasswordHash string) error {
	query := `UPDATE user SET password_hash = ?, updated_at = ? WHERE user_id = ?`
	_, err := db.Exec(query, newPasswordHash, time.Now(), userID)
	return err
}
