package repositories

import (
	"MicroGreens/internal/models"
	"context"
	"database/sql"
	"encoding/json"
)

type MicrogreenRepository struct {
	Db *sql.DB
}

func (r *MicrogreenRepository) Create(ctx context.Context, m models.Microgreen) (models.Microgreen, error) {
	substrate, _ := json.Marshal(m.Substrate)
	tips, _ := json.Marshal(m.Tips)
	query := `INSERT INTO microgreens (name, latin_name, germination_days, harvest_days, optimal_temp, light_requirements, humidity_level, substrate, watering, growth_notes, tips, image_url, is_popular) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := r.Db.ExecContext(ctx, query, m.Name, m.LatinName, m.GerminationDays, m.HarvestDays, m.OptimalTemp, m.LightRequirements, m.HumidityLevel, substrate, m.Watering, m.GrowthNotes, tips, m.ImageURL, m.IsPopular)
	if err != nil {
		return m, err
	}
	id, _ := res.LastInsertId()
	m.ID = int(id)
	return m, nil
}

func (r *MicrogreenRepository) GetAll(ctx context.Context) ([]models.Microgreen, error) {
	rows, err := r.Db.QueryContext(ctx, `SELECT id, name, latin_name, germination_days, harvest_days, optimal_temp, light_requirements, humidity_level, substrate, watering, growth_notes, tips, image_url, is_popular FROM microgreens`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Microgreen
	for rows.Next() {
		var m models.Microgreen
		var substrateJSON, tipsJSON []byte
		if err := rows.Scan(&m.ID, &m.Name, &m.LatinName, &m.GerminationDays, &m.HarvestDays, &m.OptimalTemp, &m.LightRequirements, &m.HumidityLevel, &substrateJSON, &m.Watering, &m.GrowthNotes, &tipsJSON, &m.ImageURL, &m.IsPopular); err != nil {
			return nil, err
		}
		json.Unmarshal(substrateJSON, &m.Substrate)
		json.Unmarshal(tipsJSON, &m.Tips)
		result = append(result, m)
	}
	return result, nil
}

func (r *MicrogreenRepository) GetByID(ctx context.Context, id int) (models.Microgreen, error) {
	query := `SELECT id, name, latin_name, germination_days, harvest_days, optimal_temp, light_requirements, humidity_level, substrate, watering, growth_notes, tips, image_url, is_popular FROM microgreens WHERE id = ?`
	var m models.Microgreen
	var substrateJSON, tipsJSON []byte
	err := r.Db.QueryRowContext(ctx, query, id).Scan(&m.ID, &m.Name, &m.LatinName, &m.GerminationDays, &m.HarvestDays, &m.OptimalTemp, &m.LightRequirements, &m.HumidityLevel, &substrateJSON, &m.Watering, &m.GrowthNotes, &tipsJSON, &m.ImageURL, &m.IsPopular)
	if err != nil {
		return m, err
	}
	json.Unmarshal(substrateJSON, &m.Substrate)
	json.Unmarshal(tipsJSON, &m.Tips)
	return m, nil
}

func (r *MicrogreenRepository) Update(ctx context.Context, m models.Microgreen) error {
	substrate, _ := json.Marshal(m.Substrate)
	tips, _ := json.Marshal(m.Tips)
	query := `UPDATE microgreens SET name=?, latin_name=?, germination_days=?, harvest_days=?, optimal_temp=?, light_requirements=?, humidity_level=?, substrate=?, watering=?, growth_notes=?, tips=?, image_url=?, is_popular=? WHERE id=?`
	_, err := r.Db.ExecContext(ctx, query, m.Name, m.LatinName, m.GerminationDays, m.HarvestDays, m.OptimalTemp, m.LightRequirements, m.HumidityLevel, substrate, m.Watering, m.GrowthNotes, tips, m.ImageURL, m.IsPopular, m.ID)
	return err
}

func (r *MicrogreenRepository) Delete(ctx context.Context, id int) error {
	_, err := r.Db.ExecContext(ctx, `DELETE FROM microgreens WHERE id = ?`, id)
	return err
}
