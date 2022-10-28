package domain

import "time"

type RefreshToken struct {
	Token     string    `json:"token" bson:"token"`
	ExpiresAt time.Time `json:"expires_at" bson:"expires_at"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenInput struct {
	Token string `json:"token"`
}
