package api

import "project/internal/model"

type APIService interface {

	//测试接口
	GetDatas(req model.GetDatasReq) *model.GetDatasResp
}
