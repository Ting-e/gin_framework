package repository

import (
	"project/examples/gorm_crud/model"
	"project/pkg/logger"

	"gorm.io/gorm"
)

func GetList(gorm *gorm.DB, req model.GetListReq) ([]model.User, int, error) {
	var total int64
	var datas []model.User

	query := gorm.Model(&model.User{})

	// 先查总数
	query.Count(&total)

	// 再查分页数据
	err := query.Limit(*req.Limit).Offset(*req.Skip).Find(&datas).Error
	if err != nil {
		logger.Sugar.Error(err)
		return nil, -1, err
	}

	return datas, int(total), nil
}
