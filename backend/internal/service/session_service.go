package service

import (
	"time"

	"go_sleep_admin/internal/dto"
	"go_sleep_admin/internal/models"
	"go_sleep_admin/internal/repository"
)

type SessionService struct {
	repo *repository.SessionRepository
}

func NewSessionService(repo *repository.SessionRepository) *SessionService {
	return &SessionService{repo: repo}
}

func (s *SessionService) RecordLogin(username, ip, ua, failure string, success bool) error {
	return s.repo.CreateLoginLog(&models.LoginLog{
		Username:     username,
		Success:      success,
		IP:           ip,
		UserAgent:    ua,
		FailureCause: failure,
	})
}

func (s *SessionService) LoginLogs() ([]models.LoginLog, error) {
	return s.repo.ListLoginLogs()
}

func (s *SessionService) SaveOnlineSession(adminID int64, username, accessToken, refreshToken, ip, ua string, expiredAt time.Time) error {
	return s.repo.UpsertOnlineSession(&models.OnlineSession{
		AdminID:      adminID,
		Username:     username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		IP:           ip,
		UserAgent:    ua,
		LastActiveAt: time.Now(),
		ExpiredAt:    expiredAt,
	})
}

func (s *SessionService) Touch(accessToken string) error {
	return s.repo.Touch(accessToken, time.Now())
}

func (s *SessionService) OnlineSessions() ([]dto.OnlineSessionItem, error) {
	items, err := s.repo.ListOnlineSessions(time.Now())
	if err != nil {
		return nil, err
	}
	out := make([]dto.OnlineSessionItem, 0, len(items))
	for _, item := range items {
		out = append(out, dto.OnlineSessionItem{
			ID:           item.ID,
			AdminID:      item.AdminID,
			Username:     item.Username,
			IP:           item.IP,
			UserAgent:    item.UserAgent,
			LastActiveAt: item.LastActiveAt.Format(time.DateTime),
			ExpiredAt:    item.ExpiredAt.Format(time.DateTime),
		})
	}
	return out, nil
}

func (s *SessionService) ForceOffline(id int64) error {
	return s.repo.ForceOffline(id, time.Now())
}

func (s *SessionService) IsForcedOffline(accessToken string) (bool, error) {
	return s.repo.IsForcedOffline(accessToken)
}
