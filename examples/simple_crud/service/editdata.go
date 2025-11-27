package service

import (
	"project/examples/simple_crud/model"
	"project/examples/simple_crud/repository"
	"project/pkg/database"
)

func (s *Service) EditData(req model.EditDataReq) (res model.EditDataResp) {
	res.Code = 500
	res.Mess = "内部服务出错"
	db := database.GetMysql().GetDb()

	err := repository.EditData(db, req)
	if err != nil {
		res.Code = 400
		res.Mess = "删除失败"
		return res
	}

	res.Code = 200
	res.Mess = "成功"
	return res

}
