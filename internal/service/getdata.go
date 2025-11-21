package service

import (
	"project/internal/model"
	"project/internal/repository"
	mysqldb "project/pkg/database"
)

func (s *Service) GetData(ID string) (res model.GetDataResp) {
	res.Code = 500
	res.Mess = "内部服务出错"
	db := mysqldb.GetMysql().GetDb()

	data, err := repository.GetData(db, ID)
	if err != nil {
		res.Code = 400
		res.Mess = "获取失败"
		return res
	}

	if data == nil {
		res.Code = 204
		res.Mess = "暂无数据"
		return res
	}

	res.Code = 200
	res.Mess = "成功"
	res.Datas = *data
	return res

}
