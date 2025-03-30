package repositories

import (
	"MicroGreens/internal/models"
	"context"
	"database/sql"
)

type PhotoRepository struct {
	Db *sql.DB
}

func (r *PhotoRepository) Create(ctx context.Context, p models.ObservationPhoto) (models.ObservationPhoto, error) {
	query := `INSERT INTO observation_photos (observation_id, photo_url, label) VALUES (?, ?, ?)`
	res, err := r.Db.ExecContext(ctx, query, p.ObservationID, p.PhotoURL, p.Label)
	if err != nil {
		return p, err
	}
	id, _ := res.LastInsertId()
	p.ID = int(id)
	return p, nil
}

func (r *PhotoRepository) GetAll(ctx context.Context) ([]models.ObservationPhoto, error) {
	rows, err := r.Db.QueryContext(ctx, `SELECT id, observation_id, photo_url, label, created_at FROM observation_photos`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.ObservationPhoto
	for rows.Next() {
		var p models.ObservationPhoto
		if err := rows.Scan(&p.ID, &p.ObservationID, &p.PhotoURL, &p.Label, &p.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, nil
}

func (r *PhotoRepository) GetByID(ctx context.Context, id int) (models.ObservationPhoto, error) {
	var p models.ObservationPhoto
	err := r.Db.QueryRowContext(ctx, `SELECT id, observation_id, photo_url, label, created_at FROM observation_photos WHERE id = ?`, id).
		Scan(&p.ID, &p.ObservationID, &p.PhotoURL, &p.Label, &p.CreatedAt)
	return p, err
}

func (r *PhotoRepository) Update(ctx context.Context, p models.ObservationPhoto) error {
	query := `UPDATE observation_photos SET observation_id=?, photo_url=?, label=? WHERE id=?`
	_, err := r.Db.ExecContext(ctx, query, p.ObservationID, p.PhotoURL, p.Label, p.ID)
	return err
}

func (r *PhotoRepository) Delete(ctx context.Context, id int) error {
	_, err := r.Db.ExecContext(ctx, `DELETE FROM observation_photos WHERE id = ?`, id)
	return err
}
