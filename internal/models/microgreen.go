package models

type Microgreen struct {
	ID                int      `json:"id"`
	Name              string   `json:"name"`
	LatinName         string   `json:"latin_name"`
	GerminationDays   int      `json:"germination_days"`
	HarvestDays       int      `json:"harvest_days"`
	OptimalTemp       string   `json:"optimal_temp"`
	LightRequirements string   `json:"light_requirements"`
	HumidityLevel     string   `json:"humidity_level"`
	Substrate         []string `json:"substrate"`
	Watering          string   `json:"watering"`
	GrowthNotes       string   `json:"growth_notes"`
	Tips              []string `json:"tips"`
	ImageURL          string   `json:"image_url"`
	IsPopular         bool     `json:"is_popular"`
}
