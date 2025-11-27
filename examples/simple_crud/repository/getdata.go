package repository

import (
	"database/sql"
	"project/examples/simple_crud/model"
	"project/pkg/logger"
)

func GetData(db *sql.DB, ID string) (*model.Data, error) {
	var data model.Data

	query := `
	SELECT 
		DISTINCT(key),
		COALESCE(field,""),
		DATE_FORMAT(create_time,'%Y-%m-%d %H:%i:%s')
	FROM 
		ceshi
	WHERE 
    	key = ?;`

	stmt, err := db.Prepare(query)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(ID).Scan(
		&data.Key,
		&data.Field,
		&data.Create_time,
	)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}

	return &data, nil
}
