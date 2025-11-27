package repository

import (
	"database/sql"
	"project/examples/simeple_crud/model"
	"project/pkg/logger"
	"strings"
)

func GetList(db *sql.DB, req model.GetListReq) ([]model.Data, int, error) {
	var datas []model.Data
	var total int

	query := `
	SELECT 
		DISTINCT(key),
		COALESCE(field,""),
		DATE_FORMAT(create_time,'%Y-%m-%d %H:%i:%s')
	FROM 
		ceshi
	WHERE 
    	1=1`
	count := `SELECT COUNT(DISTINCT(key)) FROM ceshi WHERE 1=1 `

	var queryArgs []interface{}
	var totalArgs []interface{}
	var queryBuilder strings.Builder
	var totalBuilder strings.Builder
	queryBuilder.WriteString(query)
	totalBuilder.WriteString(count)

	if req.Limit != nil && req.Skip != nil {
		queryArgs = append(queryArgs, req.Limit, req.Skip)
		queryBuilder.WriteString(` ORDER BY create_time DESC LIMIT ?, ?`)
	}

	stmt, err := db.Prepare(queryBuilder.String())
	if err != nil {
		logger.Sugar.Error(err)
		return nil, -1, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(queryArgs...)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, -1, err
	}
	defer rows.Close()

	for rows.Next() {
		var data model.Data
		err := rows.Scan(
			&data.Key,
			&data.Field,
			&data.Create_time,
		)
		if err != nil {
			logger.Sugar.Error(err)
			return nil, -1, err
		}

		datas = append(datas, data)
	}

	err = db.QueryRow(totalBuilder.String(), totalArgs...).Scan(&total)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, -1, err
	}

	return datas, total, nil
}
