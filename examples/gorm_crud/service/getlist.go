package service

import (
	"project/examples/gorm_crud/model"
	"project/examples/gorm_crud/repository"
	"project/pkg/database"
)

func (s *Service) GetList(req model.GetListReq) (res model.GetListResp) {
	res.Code = 500
	res.Mess = "内部服务出错"
	gormdb := database.GetGormMysql().GetDb()

	datas, total, err := repository.GetList(gormdb, req)
	if err != nil {
		res.Code = 400
		res.Mess = "获取失败"
		return res
	}

	if len(datas) == 0 {
		res.Code = 204
		res.Mess = "暂无数据"
		return res
	}

	res.Code = 200
	res.Mess = "成功"
	res.Total = total
	res.Datas = datas
	return res

}
