package service

import (
	"go_sleep_admin/internal/dto"
	"go_sleep_admin/internal/models"
	"go_sleep_admin/internal/platform/auth"
	"go_sleep_admin/internal/platform/storage"
	"go_sleep_admin/internal/repository"
)

type UploadService struct {
	repo          *repository.UploadRepository
	adminRepo     *repository.AdminRepository
	dataScopeRepo *repository.DataScopeRepository
}

func NewUploadService(repo *repository.UploadRepository, adminRepo *repository.AdminRepository, dataScopeRepo *repository.DataScopeRepository) *UploadService {
	return &UploadService{repo: repo, adminRepo: adminRepo, dataScopeRepo: dataScopeRepo}
}

func (s *UploadService) Upload(adminID int64, uploaded *storage.UploadedFile) error {
	record := &models.UploadFile{
		Type:       uploaded.Type,
		OriginName: uploaded.OriginName,
		Random:     uploaded.Random,
		Path:       uploaded.Path,
		Md5:        uploaded.MD5,
		Scene:      uploaded.Scene,
		Provider:   uploaded.Provider,
		AdminID:    adminID,
	}
	return s.repo.Save(record)
}

func (s *UploadService) List(query dto.UploadListQuery, claims *auth.Claims) ([]models.UploadFile, int64, error) {
	scope, err := s.resolveScope(claims)
	if err != nil {
		return nil, 0, err
	}
	return s.repo.List(query, scope)
}

func (s *UploadService) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *UploadService) DeleteBatch(ids []int64) error {
	return s.repo.DeleteBatch(ids)
}

func (s *UploadService) resolveScope(claims *auth.Claims) (repository.UploadScope, error) {
	if claims == nil {
		return repository.UploadScope{}, nil
	}
	admin, err := s.adminRepo.ByID(claims.AdminID)
	if err != nil {
		return repository.UploadScope{}, err
	}
	scope, err := resolveAdminScope(admin, s.dataScopeRepo)
	if err != nil {
		return repository.UploadScope{}, err
	}
	return repository.UploadScope{
		AllData: scope.AllData,
		AdminID: scope.AdminID,
		DeptIDs: scope.DeptIDs,
	}, nil
}
