package dao

import (
	"database/sql"
	"project/basic/logger"
)

func CheckWhetherAction(db *sql.DB, table string, field string, field_value string) (bool, error) {
	var num int
	query := `SELECT
	COUNT(*)
  FROM
    ` + table + `
  WHERE
     ` + field + ` = ? `
	// 检查数据是否已存在
	err := db.QueryRow(query, field_value).Scan(
		&num,
	)
	if err != nil {
		logger.Sugar.Error(err)
		return false, err
	}
	if num != 0 {
		return false, nil
	}
	return true, nil
}
