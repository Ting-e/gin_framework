package dao

import (
	"database/sql"
	"project/basic/logger"
	"project/web/handler/model"
	"strings"
	"time"
)

// 新增
func AddData(db *sql.DB, key, field string) error {
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
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		key,
		field,
		time.Now().Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	return nil
}

// 编辑
func EditData(db *sql.DB, field, key string) error {
	query := `UPDATE ceshi SET field = ? WHERE key = ? `
	_, err := db.Exec(query,
		field,
		key,
	)

	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	return err
}

// 不包含子查询
func GetData(db *sql.DB, req model.GetDatasReq) ([]*model.Data, int, error) {
	var datas []*model.Data
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

	if req.Key != "" {
		queryArgs = append(queryArgs, req.Key)
		totalArgs = append(totalArgs, req.Key)
		queryBuilder.WriteString(` AND key = ? `)
		totalBuilder.WriteString(` AND key = ?  `)
	}

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
		data := new(model.Data)
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

// 删除
func DelData(db *sql.DB, key string) error {
	query := `DELETE FROM ceshi WHERE key = ?`
	_, err := db.Exec(query, key)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	return err
}
