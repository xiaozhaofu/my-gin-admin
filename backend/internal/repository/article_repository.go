package repository

import (
	"time"

	"go_sleep_admin/internal/dto"
	"go_sleep_admin/internal/models"
	"gorm.io/gorm"
)

type ArticleRepository struct {
	db *gorm.DB
}

type ArticleScope struct {
	AllData bool
	AdminID int64
	DeptIDs []int64
}

func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

func (r *ArticleRepository) List(query dto.ArticleQuery, scope ArticleScope) ([]models.Article, int64, error) {
	var (
		items []models.Article
		total int64
	)
	page, pageSize := query.Normalize()
	db := r.db.Model(&models.Article{}).Preload("Content").Order("id desc")
	if !scope.AllData {
		switch {
		case scope.AdminID > 0:
			db = db.Where("articles.admin_id = ?", scope.AdminID)
		case len(scope.DeptIDs) > 0:
			db = db.Joins("JOIN admins ON admins.id = articles.admin_id").Where("admins.dept_id IN ?", scope.DeptIDs)
		}
	}
	if query.Title != "" {
		db = db.Where("title LIKE ?", "%"+query.Title+"%")
	}
	if query.Type != nil {
		db = db.Where("type = ?", *query.Type)
	}
	if query.MenuID != nil {
		db = db.Where("menu_id = ?", *query.MenuID)
	}
	if query.ChannelID != nil {
		db = db.Where("channel_id = ?", *query.ChannelID)
	}
	if query.IsPaid != nil {
		db = db.Where("is_paid = ?", *query.IsPaid)
	}
	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}
	if query.CreatedFrom != "" {
		if t, err := time.Parse(time.DateTime, query.CreatedFrom); err == nil {
			db = db.Where("created_at >= ?", t)
		}
	}
	if query.CreatedTo != "" {
		if t, err := time.Parse(time.DateTime, query.CreatedTo); err == nil {
			db = db.Where("created_at <= ?", t)
		}
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error
	return items, total, err
}

func (r *ArticleRepository) ByID(id int64) (*models.Article, error) {
	var item models.Article
	if err := r.db.Preload("Content").First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ArticleRepository) Save(article *models.Article, content string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return r.saveWithTx(tx, article, content)
	})
}

func (r *ArticleRepository) BatchCreate(items []*models.Article, contents []string) ([]int64, error) {
	ids := make([]int64, 0, len(items))
	err := r.db.Transaction(func(tx *gorm.DB) error {
		for idx, item := range items {
			if err := r.saveWithTx(tx, item, contents[idx]); err != nil {
				return err
			}
			ids = append(ids, item.ID)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func (r *ArticleRepository) saveWithTx(tx *gorm.DB, article *models.Article, content string) error {
	if err := tx.Save(article).Error; err != nil {
		return err
	}
	var detail models.ArticleContent
	err := tx.Where("article_id = ?", article.ID).First(&detail).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	detail.ArticleID = article.ID
	detail.Content = content
	return tx.Save(&detail).Error
}

func (r *ArticleRepository) Delete(ids []int64) error {
	return r.db.Where("id IN ?", ids).Delete(&models.Article{}).Error
}

func (r *ArticleRepository) UpdateStatus(ids []int64, status int) error {
	return r.db.Model(&models.Article{}).Where("id IN ?", ids).Update("status", status).Error
}
