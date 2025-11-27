package repository

import (
	"project/examples/gorm_crud/model"
	"project/pkg/logger"

	"gorm.io/gorm"
)

// 删除
func DelData(gorm *gorm.DB, ID string) error {

	err := gorm.Delete(model.User{
		User_id: ID,
	}).Error
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	return err
}
