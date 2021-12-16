package types

import "time"

type CreateUserInput struct {
	Email    string `validate:"email,required" json:"email"`
	Password string `validate:"gte=3,lte=100,required" json:"password"`
}

type UpdateEmailInput struct {
	ID       string `json:"-"`
	Email    string `validate:"email,required" json:"email"`
	Password string `validate:"gte=3,lte=100,required" json:"password"`
}

type UpdatePasswordInput struct {
	ID              string `json:"-"`
	CurrentPassword string `validate:"gte=3,lte=100,required" json:"current_password"`
	NewPassword     string `validate:"gte=3,lte=100,required" json:"new_password"`
}

type GetUserInput struct {
	ID string `json:"-"`
}

type DeleteUserInput struct {
	ID string `json:"-"`
}

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
