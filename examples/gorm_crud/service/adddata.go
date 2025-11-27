package service

import (
	"project/examples/gorm_crud/model"
	"project/examples/gorm_crud/repository"
	"project/pkg/database"
)

func (s *Service) AddData(req model.AddDataReq) (res model.AddDataResp) {
	res.Code = 500
	res.Mess = "内部服务出错"
	gormdb := database.GetGormMysql().GetDb()

	err := repository.AddData(gormdb, req)
	if err != nil {
		res.Code = 400
		res.Mess = "新增失败"
		return res
	}

	res.Code = 200
	res.Mess = "成功"
	return res

}
