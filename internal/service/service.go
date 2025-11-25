package service

import (
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
}
