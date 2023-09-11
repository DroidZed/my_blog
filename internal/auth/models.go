package auth

type LoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Jwt       string `json:"accessToken,omitempty"`
	Refresh   string `json:"refreshToken,omitempty"`
	UserId    string `json:"userId,omitempty"`
	Role      string `json:"role,omitempty"`
	AccStatus int    `json:"accStatus,omitempty"`
	Error     string `json:"error,omitempty"`
}

type JwtResponse struct {
	Jwt     string `json:"accessToken,omitempty"`
	Refresh string `json:"refreshToken,omitempty"`
	Error   string `json:"error,omitempty"`
}

type RefreshReq struct {
	Expired string `json:"expired"`
}
