package domain

type User struct {
	ID           string       `json:"id" bson:"_id,omitempty"`
	Username     string       `json:"username" bson:"username"`
	Email        string       `json:"email" bson:"email"`
	Password     string       `json:"-" bson:"password"`
	RefreshToken RefreshToken `json:"-" bson:"refreshToken"`
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
