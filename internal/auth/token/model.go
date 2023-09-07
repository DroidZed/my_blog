package token

type JwtResponse struct {
	Jwt     string `json:"accessToken,omitempty"`
	Refresh string `json:"refreshToken,omitempty"`
	Error   string `json:"error,omitempty"`
}
