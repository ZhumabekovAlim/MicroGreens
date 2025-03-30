package models

type ObservationPhoto struct {
	ID            int    `json:"id"`
	ObservationID int    `json:"observation_id"`
	PhotoURL      string `json:"photo_url"`
	Label         string `json:"label"`
	CreatedAt     string `json:"created_at"`
}
