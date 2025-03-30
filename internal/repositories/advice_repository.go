package repositories

import (
	"MicroGreens/internal/models"
	"context"
	"database/sql"
)

type AdviceRepository struct {
	Db *sql.DB
}

func (r *AdviceRepository) Create(ctx context.Context, a models.AdviceMessage) (models.AdviceMessage, error) {
	query := `INSERT INTO advice_messages (microgreen_id, message) VALUES (?, ?)`
	res, err := r.Db.ExecContext(ctx, query, a.MicrogreenID, a.Message)
	if err != nil {
		return a, err
	}
	id, _ := res.LastInsertId()
	a.ID = int(id)
	return a, nil
}

func (r *AdviceRepository) GetAll(ctx context.Context) ([]models.AdviceMessage, error) {
	rows, err := r.Db.QueryContext(ctx, `SELECT id, microgreen_id, message, created_at FROM advice_messages`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.AdviceMessage
	for rows.Next() {
		var a models.AdviceMessage
		if err := rows.Scan(&a.ID, &a.MicrogreenID, &a.Message, &a.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, a)
	}
	return result, nil
}

func (r *AdviceRepository) GetByID(ctx context.Context, id int) (models.AdviceMessage, error) {
	var a models.AdviceMessage
	err := r.Db.QueryRowContext(ctx, `SELECT id, microgreen_id, message, created_at FROM advice_messages WHERE id = ?`, id).
		Scan(&a.ID, &a.MicrogreenID, &a.Message, &a.CreatedAt)
	return a, err
}

func (r *AdviceRepository) Update(ctx context.Context, a models.AdviceMessage) error {
	query := `UPDATE advice_messages SET microgreen_id=?, message=? WHERE id=?`
	_, err := r.Db.ExecContext(ctx, query, a.MicrogreenID, a.Message, a.ID)
	return err
}

func (r *AdviceRepository) Delete(ctx context.Context, id int) error {
	_, err := r.Db.ExecContext(ctx, `DELETE FROM advice_messages WHERE id = ?`, id)
	return err
}
