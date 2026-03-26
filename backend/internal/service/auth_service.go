package service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"go_sleep_admin/internal/dto"
	"go_sleep_admin/internal/models"
	"go_sleep_admin/internal/platform/auth"
	"go_sleep_admin/internal/repository"
)

type AuthService struct {
	adminRepo  *repository.AdminRepository
	enforcer   Enforcer
	jwtManager *auth.JWTManager
	sessionSvc *SessionService
}

type Enforcer interface {
	GetPermissionsForUser(user string, domain ...string) ([][]string, error)
}

func NewAuthService(adminRepo *repository.AdminRepository, enforcer Enforcer, jwtManager *auth.JWTManager, sessionSvc *SessionService) *AuthService {
	return &AuthService{adminRepo: adminRepo, enforcer: enforcer, jwtManager: jwtManager, sessionSvc: sessionSvc}
}

func (s *AuthService) Login(req dto.LoginRequest, ip, userAgent string) (*dto.LoginResponse, error) {
	admin, err := s.adminRepo.ByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = s.sessionSvc.RecordLogin(req.Username, ip, userAgent, "账号不存在", false)
			return nil, errors.New("账号或密码错误")
		}
		return nil, err
	}
	if admin.Status != models.AdminStatusNormal {
		_ = s.sessionSvc.RecordLogin(req.Username, ip, userAgent, "账号已禁用", false)
		return nil, errors.New("账号已禁用")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)); err != nil {
		_ = s.sessionSvc.RecordLogin(req.Username, ip, userAgent, "密码错误", false)
		return nil, errors.New("账号或密码错误")
	}

	roles := make([]string, 0, len(admin.Roles))
	for _, role := range admin.Roles {
		roles = append(roles, role.Code)
	}

	accessToken, accessExp, err := s.jwtManager.GenerateAccessToken(admin.ID, admin.Username, roles)
	if err != nil {
		return nil, err
	}
	refreshToken, refreshExp, err := s.jwtManager.GenerateRefreshToken(admin.ID, admin.Username, roles)
	if err != nil {
		return nil, err
	}
	_ = s.sessionSvc.RecordLogin(req.Username, ip, userAgent, "", true)
	_ = s.sessionSvc.SaveOnlineSession(admin.ID, admin.Username, accessToken, refreshToken, ip, userAgent, accessExp)

	return &dto.LoginResponse{
		AccessToken:         accessToken,
		AccessTokenExpires:  accessExp.Unix(),
		RefreshToken:        refreshToken,
		RefreshTokenExpires: refreshExp.Unix(),
		User: dto.AdminBrief{
			ID:          admin.ID,
			Username:    admin.Username,
			Nickname:    admin.Nickname,
			Avatar:      admin.Avatar,
			Roles:       roles,
			Permissions: permissionsFor(roles, s.enforcer),
		},
	}, nil
}

func (s *AuthService) Profile(adminID int64) (*dto.AdminBrief, error) {
	admin, err := s.adminRepo.ByID(adminID)
	if err != nil {
		return nil, err
	}
	roles := make([]string, 0, len(admin.Roles))
	for _, role := range admin.Roles {
		roles = append(roles, role.Code)
	}
	return &dto.AdminBrief{
		ID:          admin.ID,
		Username:    admin.Username,
		Nickname:    admin.Nickname,
		Avatar:      admin.Avatar,
		Roles:       roles,
		Permissions: permissionsFor(roles, s.enforcer),
	}, nil
}

func permissionsFor(roles []string, enforcer Enforcer) []string {
	set := map[string]struct{}{}
	for _, role := range roles {
		perms, err := enforcer.GetPermissionsForUser(role)
		if err != nil {
			continue
		}
		for _, perm := range perms {
			if len(perm) > 2 {
				set[perm[1]+"#"+perm[2]] = struct{}{}
			}
		}
	}
	out := make([]string, 0, len(set))
	for item := range set {
		out = append(out, item)
	}
	return out
}
