package repositories

import (
	"MicroGreens/internal/models"
	"context"
	"database/sql"
)

type ObservationRepository struct {
	Db *sql.DB
}

func (r *ObservationRepository) Create(ctx context.Context, o models.Observation) (models.Observation, error) {
	query := `INSERT INTO observations (batch_id, date, note, height_cm, water_status, light_type, humidity_percent) VALUES (?, ?, ?, ?, ?, ?, ?)`
	res, err := r.Db.ExecContext(ctx, query, o.BatchID, o.Date, o.Note, o.HeightCM, o.WaterStatus, o.LightType, o.HumidityPercent)
	if err != nil {
		return o, err
	}
	id, _ := res.LastInsertId()
	o.ID = int(id)
	return o, nil
}

func (r *ObservationRepository) GetAll(ctx context.Context) ([]models.Observation, error) {
	rows, err := r.Db.QueryContext(ctx, `SELECT id, batch_id, date, note, height_cm, water_status, light_type, humidity_percent, created_at FROM observations`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.Observation
	for rows.Next() {
		var o models.Observation
		if err := rows.Scan(&o.ID, &o.BatchID, &o.Date, &o.Note, &o.HeightCM, &o.WaterStatus, &o.LightType, &o.HumidityPercent, &o.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, o)
	}
	return result, nil
}

func (r *ObservationRepository) GetByID(ctx context.Context, id int) (models.Observation, error) {
	var o models.Observation
	err := r.Db.QueryRowContext(ctx, `SELECT id, batch_id, date, note, height_cm, water_status, light_type, humidity_percent, created_at FROM observations WHERE id = ?`, id).
		Scan(&o.ID, &o.BatchID, &o.Date, &o.Note, &o.HeightCM, &o.WaterStatus, &o.LightType, &o.HumidityPercent, &o.CreatedAt)
	return o, err
}

func (r *ObservationRepository) Update(ctx context.Context, o models.Observation) error {
	query := `UPDATE observations SET batch_id=?, date=?, note=?, height_cm=?, water_status=?, light_type=?, humidity_percent=? WHERE id=?`
	_, err := r.Db.ExecContext(ctx, query, o.BatchID, o.Date, o.Note, o.HeightCM, o.WaterStatus, o.LightType, o.HumidityPercent, o.ID)
	return err
}

func (r *ObservationRepository) Delete(ctx context.Context, id int) error {
	_, err := r.Db.ExecContext(ctx, `DELETE FROM observations WHERE id = ?`, id)
	return err
}
