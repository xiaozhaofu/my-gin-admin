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
	menuIDs, err := resolveArticleMenuIDs(req.MenuID, req.MenuIDs)
	if err != nil {
		return err
	}
	coverLarge, coverMedium, coverSmall := resolveArticleCovers(req.Type, req.Cover, req.CoverLarge, req.CoverMedium, req.CoverSmall)
	if err := validateArticleCovers(req.Type, coverLarge, coverMedium, coverSmall); err != nil {
		return err
	}
	item := &models.Article{
		BaseID:      models.BaseID{ID: id},
		Title:       req.Title,
		Summary:     req.Summary,
		Type:        models.ArticleType(req.Type),
		CoverLarge:  coverLarge,
		CoverMedium: coverMedium,
		CoverSmall:  coverSmall,
		CoverType:   req.CoverType,
		MenuID:      menuIDs[0],
		ChannelID:   req.ChannelID,
		SortOrder:   req.SortOrder,
		IsPaid:      models.BooleanField(req.IsPaid),
		AdminID:     adminID,
		IsTop:       models.BooleanField(req.IsTop),
		IsHot:       models.BooleanField(req.IsHot),
		IsRecommend: models.BooleanField(req.IsRecommend),
		Status:      models.ArticleStatus(req.Status),
	}
	return s.repo.Save(item, req.Content, menuIDs)
}

func (s *ArticleService) BatchCreate(adminID int64, req dto.ArticleBatchCreateRequest) ([]int64, error) {
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("至少需要一条文章数据")
	}
	if len(req.Items) > 100 {
		return nil, fmt.Errorf("单次最多批量新增 100 篇文章")
	}
	menuIDs, err := resolveArticleMenuIDs(req.MenuID, req.MenuIDs)
	if err != nil {
		return nil, err
	}

	articles := make([]*models.Article, 0, len(req.Items))
	contents := make([]string, 0, len(req.Items))
	menuIDsByArticle := make([][]int64, 0, len(req.Items))
	for idx, row := range req.Items {
		title := strings.TrimSpace(row.Title)
		content := strings.TrimSpace(row.Content)
		if title == "" {
			return nil, fmt.Errorf("第 %d 行标题不能为空", idx+1)
		}
		if content == "" {
			return nil, fmt.Errorf("第 %d 行正文不能为空", idx+1)
		}
		coverLarge, coverMedium, coverSmall := resolveArticleCovers(
			req.Type,
			firstNonEmpty(row.Cover, req.Cover),
			firstNonEmpty(row.CoverLarge, req.CoverLarge),
			firstNonEmpty(row.CoverMedium, req.CoverMedium),
			firstNonEmpty(row.CoverSmall, req.CoverSmall),
		)
		if err := validateArticleCovers(req.Type, coverLarge, coverMedium, coverSmall); err != nil {
			return nil, fmt.Errorf("第 %d 行%s", idx+1, err.Error())
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
			CoverLarge:  coverLarge,
			CoverMedium: coverMedium,
			CoverSmall:  coverSmall,
			CoverType:   coverType,
			MenuID:      menuIDs[0],
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
		menuIDsByArticle = append(menuIDsByArticle, menuIDs)
	}

	return s.repo.BatchCreate(articles, contents, menuIDsByArticle)
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

func resolveArticleMenuIDs(primary int64, menuIDs []int64) ([]int64, error) {
	resolved := make([]int64, 0, len(menuIDs)+1)
	seen := make(map[int64]struct{}, len(menuIDs)+1)
	appendMenu := func(menuID int64) {
		if menuID <= 0 {
			return
		}
		if _, ok := seen[menuID]; ok {
			return
		}
		seen[menuID] = struct{}{}
		resolved = append(resolved, menuID)
	}

	for _, menuID := range menuIDs {
		appendMenu(menuID)
	}
	appendMenu(primary)
	if len(resolved) == 0 {
		return nil, fmt.Errorf("请至少选择一个内容菜单")
	}
	return resolved, nil
}

func resolveArticleCovers(articleType int, cover, coverLarge, coverMedium, coverSmall string) (string, string, string) {
	large := strings.TrimSpace(firstNonEmpty(coverLarge, cover))
	medium := strings.TrimSpace(firstNonEmpty(coverMedium, cover))
	small := strings.TrimSpace(firstNonEmpty(coverSmall, cover))

	if articleType != int(models.ArticleTypeAudio) {
		return large, "", ""
	}
	return large, medium, small
}

func validateArticleCovers(articleType int, coverLarge, coverMedium, coverSmall string) error {
	if strings.TrimSpace(coverLarge) == "" {
		return fmt.Errorf("缺少大封面图")
	}
	if articleType == int(models.ArticleTypeAudio) {
		if strings.TrimSpace(coverMedium) == "" || strings.TrimSpace(coverSmall) == "" {
			return fmt.Errorf("音频文章需要同时上传中封面图和小封面图")
		}
	}
	return nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
