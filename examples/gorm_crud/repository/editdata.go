package repository

import (
	"log"
	"project/examples/gorm_crud/model"

	"gorm.io/gorm"
)

// 编辑
func EditData(gorm *gorm.DB, req model.EditDataReq) error {
	user := model.User{
		User_id: req.ID,
		Name:    req.Name,
	}

	err := gorm.Model(user).Update("name", user.Name).Where("user_id=?", user).Error
	if err != nil {
		log.Print("数据修改失败", err)
	}

	return err
}
