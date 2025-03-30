package models

type Observation struct {
	ID              int     `json:"id"`
	BatchID         int     `json:"batch_id"`
	Date            string  `json:"date"`
	Note            string  `json:"note"`
	HeightCM        float64 `json:"height_cm"`
	WaterStatus     string  `json:"water_status"`
	LightType       string  `json:"light_type"`
	HumidityPercent int     `json:"humidity_percent"`
	CreatedAt       string  `json:"created_at"`
}
