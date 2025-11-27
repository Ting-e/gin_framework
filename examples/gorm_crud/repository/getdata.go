package repository

import (
	"project/examples/gorm_crud/model"
	"project/pkg/logger"

	"gorm.io/gorm"
)

func GetData(gorm *gorm.DB, ID string) (*model.User, error) {
	var data model.User

	err := gorm.Find(&data).Where("user_id=?", ID).Error

	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}

	return &data, nil
}
