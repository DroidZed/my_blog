package auth

type LoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Jwt     string `json:"accessToken,omitempty"`
	Refresh string `json:"refreshToken,omitempty"`
	Error   string `json:"error,omitempty"`
}
