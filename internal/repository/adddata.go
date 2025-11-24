package repository

import (
	"database/sql"
	"project/internal/model"
	"project/pkg/logger"
	"project/pkg/utils/idgen"
	"time"
)

// 新增
func AddData(db *sql.DB, req model.AddDataReq) error {

	key, err := idgen.GenerateID(16)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	query := `
	INSERT INTO
		ceshi 
	(
	    key,
		field,
		create_time
	)
  	VALUES
		(?, ?, ?);`

	stmt, err := db.Prepare(query)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		key,
		req.Field,
		time.Now().Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	return nil
}
