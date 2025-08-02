package models

type User struct {
	UID   string `json:"uid"`
	Name  string `json:"name"`
	Age   int    `json:"age" validate:"gte=16"`
	Email string `json:"email" validate:"email"`
	Pass  string `json:"pass" validate:"min=8"`
}

type UserLogin struct {
	Email string `json:"email" validate:"email"`
	Pass  string `json:"pass" validate:"min=8"`
}
