package domain

import "time"

type CreateUser struct {
	Username string
	Email    string
	Password string
	Name     string
	Surname  string
	IsMale   bool
}

type CreateUserRow struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	ID        int64
	Username  string
	Email     string
	Password  string
	Name      string
	Surname   string
	IsMale    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(id int64, username, email, password, name, surname string, isMale bool,
	createdAt, updatedAt time.Time,
) User {
	return User{
		ID:        id,
		Username:  username,
		Email:     email,
		Password:  password,
		Name:      name,
		Surname:   surname,
		IsMale:    isMale,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

// DTO для обновления юзера
type UpdateUser struct {
	ID       int64
	Username *string
	Email    *string
	Name     *string
	Surname  *string
	IsMale   *bool
}

type UserCacheProfile struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	IsMale   bool   `json:"is_male"`
}
