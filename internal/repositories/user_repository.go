package repositories

import (
	"MicroGreens/internal/models"
	"context"
	"database/sql"
)

type UserRepository struct {
	Db *sql.DB
}

func (r *UserRepository) Create(ctx context.Context, u models.User) (models.User, error) {
	query := `INSERT INTO users (email, password_hash) VALUES (?, ?)`
	res, err := r.Db.ExecContext(ctx, query, u.Email, u.PasswordHash)
	if err != nil {
		return u, err
	}
	id, _ := res.LastInsertId()
	u.ID = int(id)
	return u, nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	rows, err := r.Db.QueryContext(ctx, `SELECT id, email, password_hash, created_at FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, u)
	}
	return result, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (models.User, error) {
	var u models.User
	err := r.Db.QueryRowContext(ctx, `SELECT id, email, password_hash, created_at FROM users WHERE id = ?`, id).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	return u, err
}

func (r *UserRepository) Update(ctx context.Context, u models.User) error {
	query := `UPDATE users SET email = ?, password_hash = ? WHERE id = ?`
	_, err := r.Db.ExecContext(ctx, query, u.Email, u.PasswordHash, u.ID)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id int) error {
	_, err := r.Db.ExecContext(ctx, `DELETE FROM users WHERE id = ?`, id)
	return err
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.Db.QueryRow("SELECT id, email, password_hash FROM users WHERE email = ?", email).
		Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
