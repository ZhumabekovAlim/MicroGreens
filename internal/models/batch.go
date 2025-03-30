package models

type Batch struct {
	ID                  int    `json:"id"`
	UserID              int    `json:"user_id"`
	Name                string `json:"name"`
	MicrogreenID        int    `json:"microgreen_id"`
	SowingDate          string `json:"sowing_date"`
	Substrate           string `json:"substrate"`
	Comment             string `json:"comment"`
	EstimatedHarvestDay int    `json:"estimated_harvest_days"`
	CreatedAt           string `json:"created_at"`
}
