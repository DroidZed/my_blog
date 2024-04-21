package auth

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginBody struct {
	Payload *Login `in:"body=json"`
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

type Refresh struct {
	Expired string `json:"expired"`
}

type RefreshReq struct {
	Payload *Refresh
}
