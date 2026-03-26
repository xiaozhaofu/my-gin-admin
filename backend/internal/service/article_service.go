package service

import (
	"fmt"
	"strings"

	"go_sleep_admin/internal/dto"
	"go_sleep_admin/internal/models"
	"go_sleep_admin/internal/platform/auth"
	"go_sleep_admin/internal/repository"
)

type ArticleService struct {
	repo          *repository.ArticleRepository
	adminRepo     *repository.AdminRepository
	dataScopeRepo *repository.DataScopeRepository
}

func NewArticleService(repo *repository.ArticleRepository, adminRepo *repository.AdminRepository, dataScopeRepo *repository.DataScopeRepository) *ArticleService {
	return &ArticleService{repo: repo, adminRepo: adminRepo, dataScopeRepo: dataScopeRepo}
}

func (s *ArticleService) List(query dto.ArticleQuery, claims *auth.Claims) ([]models.Article, int64, error) {
	scope, err := s.resolveScope(claims)
	if err != nil {
		return nil, 0, err
	}
	return s.repo.List(query, scope)
}

func (s *ArticleService) Detail(id int64) (*models.Article, error) {
	return s.repo.ByID(id)
}

func (s *ArticleService) Save(id int64, adminID int64, req dto.ArticleUpsertRequest) error {
	item := &models.Article{
		BaseID:      models.BaseID{ID: id},
		Title:       req.Title,
		Summary:     req.Summary,
		Type:        models.ArticleType(req.Type),
		Cover:       req.Cover,
		CoverType:   req.CoverType,
		MenuID:      req.MenuID,
		ChannelID:   req.ChannelID,
		SortOrder:   req.SortOrder,
		IsPaid:      models.BooleanField(req.IsPaid),
		AdminID:     adminID,
		IsTop:       models.BooleanField(req.IsTop),
		IsHot:       models.BooleanField(req.IsHot),
		IsRecommend: models.BooleanField(req.IsRecommend),
		Status:      models.ArticleStatus(req.Status),
	}
	return s.repo.Save(item, req.Content)
}

func (s *ArticleService) BatchCreate(adminID int64, req dto.ArticleBatchCreateRequest) ([]int64, error) {
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("至少需要一条文章数据")
	}
	if len(req.Items) > 100 {
		return nil, fmt.Errorf("单次最多批量新增 100 篇文章")
	}

	articles := make([]*models.Article, 0, len(req.Items))
	contents := make([]string, 0, len(req.Items))
	for idx, row := range req.Items {
		title := strings.TrimSpace(row.Title)
		content := strings.TrimSpace(row.Content)
		if title == "" {
			return nil, fmt.Errorf("第 %d 行标题不能为空", idx+1)
		}
		if content == "" {
			return nil, fmt.Errorf("第 %d 行正文不能为空", idx+1)
		}
		cover := strings.TrimSpace(row.Cover)
		if cover == "" {
			cover = strings.TrimSpace(req.Cover)
		}
		if cover == "" {
			return nil, fmt.Errorf("第 %d 行缺少封面图", idx+1)
		}
		coverType := strings.TrimSpace(row.CoverType)
		if coverType == "" {
			coverType = strings.TrimSpace(req.CoverType)
		}
		if coverType == "" {
			coverType = "1"
		}

		articles = append(articles, &models.Article{
			Title:       title,
			Summary:     strings.TrimSpace(row.Summary),
			Type:        models.ArticleType(req.Type),
			Cover:       cover,
			CoverType:   coverType,
			MenuID:      req.MenuID,
			ChannelID:   req.ChannelID,
			SortOrder:   req.SortOrder,
			IsPaid:      models.BooleanField(req.IsPaid),
			AdminID:     adminID,
			IsTop:       models.BooleanField(req.IsTop),
			IsHot:       models.BooleanField(req.IsHot),
			IsRecommend: models.BooleanField(req.IsRecommend),
			Status:      models.ArticleStatus(req.Status),
		})
		contents = append(contents, content)
	}

	return s.repo.BatchCreate(articles, contents)
}

func (s *ArticleService) Delete(ids []int64) error { return s.repo.Delete(ids) }
func (s *ArticleService) UpdateStatus(ids []int64, status int) error {
	return s.repo.UpdateStatus(ids, status)
}

func (s *ArticleService) resolveScope(claims *auth.Claims) (repository.ArticleScope, error) {
	if claims == nil {
		return repository.ArticleScope{}, nil
	}
	admin, err := s.adminRepo.ByID(claims.AdminID)
	if err != nil {
		return repository.ArticleScope{}, err
	}
	scope, err := resolveAdminScope(admin, s.dataScopeRepo)
	if err != nil {
		return repository.ArticleScope{}, err
	}
	return repository.ArticleScope{
		AllData: scope.AllData,
		AdminID: scope.AdminID,
		DeptIDs: scope.DeptIDs,
	}, nil
}
