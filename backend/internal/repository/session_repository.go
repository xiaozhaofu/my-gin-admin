package repository

import (
	"time"

	"go_sleep_admin/internal/models"
	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository { return &SessionRepository{db: db} }

func (r *SessionRepository) CreateLoginLog(item *models.LoginLog) error {
	return r.db.Create(item).Error
}

func (r *SessionRepository) ListLoginLogs() ([]models.LoginLog, error) {
	var items []models.LoginLog
	err := r.db.Order("id desc").Limit(200).Find(&items).Error
	return items, err
}

func (r *SessionRepository) UpsertOnlineSession(item *models.OnlineSession) error {
	return r.db.Where("access_token = ?", item.AccessToken).Assign(item).FirstOrCreate(item).Error
}

func (r *SessionRepository) Touch(accessToken string, now time.Time) error {
	return r.db.Model(&models.OnlineSession{}).Where("access_token = ?", accessToken).Update("last_active_at", now).Error
}

func (r *SessionRepository) ListOnlineSessions(now time.Time) ([]models.OnlineSession, error) {
	var items []models.OnlineSession
	err := r.db.Where("expired_at > ?", now).Where("force_offline_at IS NULL").Order("last_active_at desc").Find(&items).Error
	return items, err
}

func (r *SessionRepository) ForceOffline(id int64, now time.Time) error {
	return r.db.Model(&models.OnlineSession{}).Where("id = ?", id).Update("force_offline_at", now).Error
}

func (r *SessionRepository) IsForcedOffline(accessToken string) (bool, error) {
	var count int64
	err := r.db.Model(&models.OnlineSession{}).
		Where("access_token = ?", accessToken).
		Where("force_offline_at IS NOT NULL").
		Count(&count).Error
	return count > 0, err
}
