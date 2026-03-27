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
		db = db.Where(
			"articles.menu_id = ? OR EXISTS (SELECT 1 FROM article_menus am WHERE am.article_id = articles.id AND am.menu_id = ?)",
			*query.MenuID,
			*query.MenuID,
		)
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
	if err := db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	if err := r.populateMenuIDs(&items); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *ArticleRepository) ByID(id int64) (*models.Article, error) {
	var item models.Article
	if err := r.db.Preload("Content").First(&item, id).Error; err != nil {
		return nil, err
	}
	if err := r.populateMenuIDs(&item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ArticleRepository) Save(article *models.Article, content string, menuIDs []int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return r.saveWithTx(tx, article, content, menuIDs)
	})
}

func (r *ArticleRepository) BatchCreate(items []*models.Article, contents []string, menuIDs [][]int64) ([]int64, error) {
	ids := make([]int64, 0, len(items))
	err := r.db.Transaction(func(tx *gorm.DB) error {
		for idx, item := range items {
			if err := r.saveWithTx(tx, item, contents[idx], menuIDs[idx]); err != nil {
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

func (r *ArticleRepository) saveWithTx(tx *gorm.DB, article *models.Article, content string, menuIDs []int64) error {
	if article.ID == 0 {
		if err := tx.Create(article).Error; err != nil {
			return err
		}
	} else {
		updates := map[string]any{
			"title":        article.Title,
			"summary":      article.Summary,
			"type":         article.Type,
			"cover_large":  article.CoverLarge,
			"cover_medium": article.CoverMedium,
			"cover_small":  article.CoverSmall,
			"cover_type":   article.CoverType,
			"menu_id":      article.MenuID,
			"channel_id":   article.ChannelID,
			"sort_order":   article.SortOrder,
			"is_paid":      article.IsPaid,
			"admin_id":     article.AdminID,
			"is_top":       article.IsTop,
			"is_hot":       article.IsHot,
			"is_recommend": article.IsRecommend,
			"status":       article.Status,
		}
		if err := tx.Model(&models.Article{}).Where("id = ?", article.ID).Updates(updates).Error; err != nil {
			return err
		}
	}

	var detail models.ArticleContent
	err := tx.Where("article_id = ?", article.ID).First(&detail).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == gorm.ErrRecordNotFound {
		detail.ArticleID = article.ID
		detail.Content = content
		if err := tx.Create(&detail).Error; err != nil {
			return err
		}
	} else {
		if err := tx.Model(&models.ArticleContent{}).Where("article_id = ?", article.ID).Update("content", content).Error; err != nil {
			return err
		}
	}
	return syncArticleMenus(tx, article.ID, menuIDs)
}

func (r *ArticleRepository) Delete(ids []int64) error {
	return r.db.Where("id IN ?", ids).Delete(&models.Article{}).Error
}

func (r *ArticleRepository) UpdateStatus(ids []int64, status int) error {
	return r.db.Model(&models.Article{}).Where("id IN ?", ids).Update("status", status).Error
}

func (r *ArticleRepository) populateMenuIDs(target any) error {
	switch items := target.(type) {
	case *models.Article:
		menuMap, err := r.fetchMenuIDs([]int64{items.ID})
		if err != nil {
			return err
		}
		items.MenuIDs = fallbackMenuIDs(items.MenuID, menuMap[items.ID])
		return nil
	case *[]models.Article:
		articleIDs := make([]int64, 0, len(*items))
		for _, item := range *items {
			articleIDs = append(articleIDs, item.ID)
		}
		menuMap, err := r.fetchMenuIDs(articleIDs)
		if err != nil {
			return err
		}
		for idx := range *items {
			(*items)[idx].MenuIDs = fallbackMenuIDs((*items)[idx].MenuID, menuMap[(*items)[idx].ID])
		}
		return nil
	default:
		return nil
	}
}

func (r *ArticleRepository) fetchMenuIDs(articleIDs []int64) (map[int64][]int64, error) {
	out := make(map[int64][]int64, len(articleIDs))
	if len(articleIDs) == 0 {
		return out, nil
	}

	var links []models.ArticleMenu
	if err := r.db.Where("article_id IN ?", articleIDs).Order("menu_id asc").Find(&links).Error; err != nil {
		return nil, err
	}
	for _, link := range links {
		out[link.ArticleID] = append(out[link.ArticleID], link.MenuID)
	}
	return out, nil
}

func syncArticleMenus(tx *gorm.DB, articleID int64, menuIDs []int64) error {
	if err := tx.Where("article_id = ?", articleID).Delete(&models.ArticleMenu{}).Error; err != nil {
		return err
	}
	uniqueMenuIDs := uniqueInt64s(menuIDs)
	if len(uniqueMenuIDs) == 0 {
		return nil
	}

	links := make([]models.ArticleMenu, 0, len(uniqueMenuIDs))
	for _, menuID := range uniqueMenuIDs {
		links = append(links, models.ArticleMenu{ArticleID: articleID, MenuID: menuID})
	}
	return tx.Create(&links).Error
}

func uniqueInt64s(items []int64) []int64 {
	out := make([]int64, 0, len(items))
	seen := make(map[int64]struct{}, len(items))
	for _, item := range items {
		if item <= 0 {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		out = append(out, item)
	}
	return out
}

func fallbackMenuIDs(primary int64, menuIDs []int64) []int64 {
	if len(menuIDs) > 0 {
		return menuIDs
	}
	if primary > 0 {
		return []int64{primary}
	}
	return nil
}
