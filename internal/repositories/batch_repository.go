package repositories

import (
	"MicroGreens/internal/models"
	"context"
	"database/sql"
)

type BatchRepository struct {
	Db *sql.DB
}

func (r *BatchRepository) Create(ctx context.Context, b models.Batch) (models.Batch, error) {
	query := `INSERT INTO batches (user_id, name, microgreen_id, sowing_date, substrate, comment, estimated_harvest_days) VALUES (?, ?, ?, ?, ?, ?, ?)`
	res, err := r.Db.ExecContext(ctx, query, b.UserID, b.Name, b.MicrogreenID, b.SowingDate, b.Substrate, b.Comment, b.EstimatedHarvestDay)
	if err != nil {
		return b, err
	}
	id, _ := res.LastInsertId()
	b.ID = int(id)
	return b, nil
}

func (r *BatchRepository) GetAll(ctx context.Context) ([]models.Batch, error) {
	rows, err := r.Db.QueryContext(ctx, `SELECT id, user_id, name, microgreen_id, sowing_date, substrate, comment, estimated_harvest_days, created_at FROM batches`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Batch
	for rows.Next() {
		var b models.Batch
		if err := rows.Scan(&b.ID, &b.UserID, &b.Name, &b.MicrogreenID, &b.SowingDate, &b.Substrate, &b.Comment, &b.EstimatedHarvestDay, &b.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, b)
	}
	return result, nil
}

func (r *BatchRepository) GetByID(ctx context.Context, id int) (models.Batch, error) {
	var b models.Batch
	err := r.Db.QueryRowContext(ctx, `SELECT id, user_id, name, microgreen_id, sowing_date, substrate, comment, estimated_harvest_days, created_at FROM batches WHERE id = ?`, id).Scan(&b.ID, &b.UserID, &b.Name, &b.MicrogreenID, &b.SowingDate, &b.Substrate, &b.Comment, &b.EstimatedHarvestDay, &b.CreatedAt)
	return b, err
}

func (r *BatchRepository) Update(ctx context.Context, b models.Batch) error {
	query := `UPDATE batches SET user_id=?, name=?, microgreen_id=?, sowing_date=?, substrate=?, comment=?, estimated_harvest_days=? WHERE id=?`
	_, err := r.Db.ExecContext(ctx, query, b.UserID, b.Name, b.MicrogreenID, b.SowingDate, b.Substrate, b.Comment, b.EstimatedHarvestDay, b.ID)
	return err
}

func (r *BatchRepository) Delete(ctx context.Context, id int) error {
	_, err := r.Db.ExecContext(ctx, `DELETE FROM batches WHERE id = ?`, id)
	return err
}
