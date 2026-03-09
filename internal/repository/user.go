package repository

import (
	"context"
	"koda-b6-backend/internal/models"
	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	db *pgx.Conn
}

func NewUserRepository(db *pgx.Conn) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	rows, err := u.db.Query(ctx,
		`SELECT id, full_name, email, password, phone FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		return nil, err
	}

	return users, nil
}


func (u *UserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	
	err := u.db.QueryRow(ctx,
		`SELECT id, full_name, email, password, phone FROM users WHERE id = $1`,
		id).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Phone)
	
	if err != nil {
		return nil, err
	}

	return &user, nil
}


func (u *UserRepository) Create(ctx context.Context, user *models.User) error {
	_, err := u.db.Exec(ctx,
		`INSERT INTO users (id, full_name, email, password, phone) 
		 VALUES ($1, $2, $3, $4, $5)`,
		user.ID, user.Name, user.Email, user.Password, user.Phone)
	
	return err
}

// Update mengubah data user
func (u *UserRepository) Update(ctx context.Context, user *models.User) error {
	_, err := u.db.Exec(ctx,
		`UPDATE users SET full_name = $1, email = $2, password = $3, phone = $4 
		 WHERE id = $5`,
		user.Name, user.Email, user.Password, user.Phone, user.ID)
	
	return err
}

// Delete menghapus user
func (u *UserRepository) Delete(ctx context.Context, id string) error {
	_, err := u.db.Exec(ctx,
		`DELETE FROM users WHERE id = $1`,
		id)
	
	return err
}