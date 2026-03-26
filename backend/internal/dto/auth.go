package dto

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LoginResponse struct {
	AccessToken         string     `json:"access_token"`
	AccessTokenExpires  int64      `json:"access_token_expires"`
	RefreshToken        string     `json:"refresh_token"`
	RefreshTokenExpires int64      `json:"refresh_token_expires"`
	User                AdminBrief `json:"user"`
}

type AdminBrief struct {
	ID          int64    `json:"id"`
	Username    string   `json:"username"`
	Nickname    string   `json:"nickname"`
	Avatar      string   `json:"avatar"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}
