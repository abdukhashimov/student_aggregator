package domain

type User struct {
	ID           string         `json:"id" bson:"_id,omitempty"`
	Username     string         `json:"username"`
	Email        string         `json:"email"`
	Password     string         `json:"-"`
	RefreshToken []RefreshToken `json:"-"`
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
