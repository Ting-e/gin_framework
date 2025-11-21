package service

import (
	"project/internal/model"
	"project/pkg/utils/snowflake"
)

var defaultService *Service

type Service struct {
	snowflake *snowflake.Worker
}

func GetService() *Service {
	return defaultService
}

func init() {
	defaultService = &Service{
		snowflake: snowflake.NewWorker(snowflake.WorkerID, snowflake.WataCenterID),
	}
}

type APIService interface {

	// 获取列表（测试接口）
	GetList(req model.GetListReq) model.GetListResp

	// 获取详情（测试接口）
	GetData(ID string) model.GetDataResp

	// 新增数据（测试接口）
	AddData(req model.AddDataReq) model.AddDataResp

	// 删除数据（测试接口）
	DelData(ID string) model.DelDataResp

	// 编辑数据（测试接口）
	EditData(req model.EditDataReq) model.EditDataResp
}
