package dto

type OnlineSessionItem struct {
	ID           int64  `json:"id"`
	AdminID      int64  `json:"admin_id"`
	Username     string `json:"username"`
	IP           string `json:"ip"`
	UserAgent    string `json:"user_agent"`
	LastActiveAt string `json:"last_active_at"`
	ExpiredAt    string `json:"expired_at"`
}
