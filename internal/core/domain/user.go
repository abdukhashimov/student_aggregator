package domain

type User struct {
	ID           string       `json:"id" bson:"_id,omitempty"`
	Username     string       `json:"username" bson:"username"`
	Email        string       `json:"email" bson:"email"`
	Password     string       `json:"-" bson:"password"`
	RefreshToken RefreshToken `json:"-" bson:"refresh_token"`
}

type UserProfile struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
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
	RefreshToken *RefreshToken
}

func (u *User) GetProfile() *UserProfile {
	return &UserProfile{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	}
}
