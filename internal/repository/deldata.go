package repository

import (
	"database/sql"
	"project/pkg/logger"
)

// 删除
func DelData(db *sql.DB, ID string) error {
	query := `DELETE FROM ceshi WHERE key = ?`
	_, err := db.Exec(query, ID)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	return err
}
