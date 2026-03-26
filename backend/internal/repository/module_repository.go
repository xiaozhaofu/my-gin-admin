package repository

import (
	"fmt"

	"gorm.io/gorm"
)

type ModuleRepository struct {
	db *gorm.DB
}

func NewModuleRepository(db *gorm.DB) *ModuleRepository { return &ModuleRepository{db: db} }

func (r *ModuleRepository) List(model interface{}, page, pageSize int) (interface{}, int64, error) {
	var total int64
	if err := r.db.Model(model).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	items := slicePtr(model)
	if err := r.db.Model(model).Order("id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *ModuleRepository) Detail(model interface{}, id int64) (interface{}, error) {
	item := valuePtr(model)
	if err := r.db.Model(model).First(item, id).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func slicePtr(model interface{}) interface{} {
	switch model.(type) {
	case *struct{}:
		panic(fmt.Sprintf("unsupported model %T", model))
	}
	return model
}

func valuePtr(model interface{}) interface{} {
	return model
}
