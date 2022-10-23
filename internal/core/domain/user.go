package domain

type User struct {
	ID           uint         `json:"id"`
	Username     string       `json:"username"`
	Email        string       `json:"email"`
	Password     string       `json:"password"`
	RefreshToken RefreshToken `json:"refresh_token"`
}

type SignUpUserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserInput struct {
	Username     *string
	Email        *string
	Password     *string
	RefreshToken *string
}
