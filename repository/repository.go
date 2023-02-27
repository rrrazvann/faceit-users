package repository

import (
	"faceit/model"

	"gorm.io/gorm"
)

type Repository[T any] struct {
	DB *gorm.DB
}

func (r *Repository[T]) Insert(model T) (*T, error) {
	result := r.DB.Create(&model)

	return &model, result.Error
}

func (r *Repository[T]) Delete(id int) error {
	var model T
	result := r.DB.Where("id = ?", id).Delete(model)

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	
	return result.Error
}

func (r *Repository[T]) Update(model T) (*T, error) {
	result := r.DB.Model(&model).Updates(&model)
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &model, result.Error
}

func (r *Repository[T]) GetAll(offset int, pageSize int, userConditions model.User) ([]T, error) {
	models := make([]T, 0)

	result := r.DB.Where(userConditions).Offset(offset).Limit(pageSize).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	return models, nil
}
