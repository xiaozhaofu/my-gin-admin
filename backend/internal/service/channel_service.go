package service

import (
	"fmt"
	"strings"

	"go_sleep_admin/internal/dto"
	"go_sleep_admin/internal/models"
	"go_sleep_admin/internal/repository"
)

type ChannelService struct {
	repo *repository.ChannelRepository
}

func NewChannelService(repo *repository.ChannelRepository) *ChannelService {
	return &ChannelService{repo: repo}
}

func (s *ChannelService) List() ([]models.Channel, error) {
	return s.repo.List()
}

func (s *ChannelService) Save(id, adminID int64, req dto.ChannelUpsertRequest) error {
	name := strings.TrimSpace(req.Name)
	code := strings.TrimSpace(req.Code)
	if name == "" {
		return fmt.Errorf("渠道名称不能为空")
	}
	if code == "" {
		return fmt.Errorf("渠道编码不能为空")
	}

	exists, err := s.repo.ExistsByNameOrCode(name, code, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("渠道名称或编码已存在")
	}

	status := models.ChannelStatus(req.Status)
	if status == 0 {
		status = models.ChannelStatusNormal
	}

	item := &models.Channel{
		ID:      id,
		Name:    name,
		Code:    code,
		Status:  status,
		Remark:  strings.TrimSpace(req.Remark),
		AdminID: adminID,
	}
	if id > 0 && adminID == 0 {
		existing, err := s.repo.ByID(id)
		if err != nil {
			return err
		}
		item.AdminID = existing.AdminID
	}
	return s.repo.Save(item)
}

func (s *ChannelService) UpdateStatus(ids []int64, status int) error {
	channelStatus := models.ChannelStatus(status)
	if channelStatus == 0 {
		return fmt.Errorf("无效的渠道状态")
	}
	return s.repo.UpdateStatus(ids, channelStatus)
}
