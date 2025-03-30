package models

type User struct {
	ID          int    `json:"id"`
	FIO         string `json:"fio"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
	CreatedAt   string `json:"created_at"`
}
