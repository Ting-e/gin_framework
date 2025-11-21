package service

import (
	"project/internal/model"
	"project/internal/repository"
	mysqldb "project/pkg/database"
)

func (s *Service) AddData(req model.AddDataReq) (res model.AddDataResp) {
	res.Code = 500
	res.Mess = "内部服务出错"
	db := mysqldb.GetMysql().GetDb()

	err := repository.AddData(db, req)
	if err != nil {
		res.Code = 400
		res.Mess = "新增失败"
		return res
	}

	res.Code = 200
	res.Mess = "成功"
	return res

}
