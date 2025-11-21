package service

import (
	mysqldb "project/basic/component/mysql"
	"project/web/handler/dao"
	"project/web/handler/model"
)

func (s *Service) GetDatas(req model.GetDatasReq) *model.GetDatasResp {
	res := new(model.GetDatasResp)
	res.Code = 500
	res.Mess = "内部服务出错"
	db := mysqldb.GetMysql().GetDb()

	datas, total, err := dao.GetData(db, req)
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
