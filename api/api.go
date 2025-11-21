package api

import "project/web/handler/model"

type APIService interface {

	//测试接口
	GetDatas(req model.GetDatasReq) *model.GetDatasResp
}
