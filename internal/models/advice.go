package models

type AdviceMessage struct {
	ID           int    `json:"id"`
	MicrogreenID int    `json:"microgreen_id"`
	Message      string `json:"message"`
	CreatedAt    string `json:"created_at"`
}
