package repository

import (
	"project/examples/gorm_crud/model"
	"project/pkg/logger"
	"project/pkg/utils/idgen"
	"time"

	"gorm.io/gorm"
)

// 新增
func AddData(gorm *gorm.DB, req model.AddDataReq) error {

	key, err := idgen.GenerateID(16)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	data := &model.User{
		User_id:     key,
		Name:        req.Name,
		Age:         req.Age,
		Create_time: time.Now(),
	}

	err = gorm.Create(data).Error
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	return nil
}
