package domain

type RefreshToken struct {
	Token     string `json:"token"`
	ExpiresAt int    `json:"expires_at"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
