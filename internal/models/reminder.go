package models

type Reminder struct {
	ID      int    `json:"id"`
	UserID  int    `json:"user_id"`
	Message string `json:"message"`
	Time    string `json:"time"`
	Active  bool   `json:"active"`
}
