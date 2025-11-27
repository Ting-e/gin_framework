package repository

import (
	"database/sql"
	"project/examples/simple_crud/model"
	"project/pkg/logger"
)

// 编辑
func EditData(db *sql.DB, req model.EditDataReq) error {
	query := `UPDATE ceshi SET field = ? WHERE key = ? `
	_, err := db.Exec(query,
		req.Field,
		req.ID,
	)

	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	return err
}
