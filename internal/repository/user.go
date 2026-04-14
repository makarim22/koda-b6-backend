package repository

import (
	"context"
	"fmt"
	"koda-b6-backend/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}


func (u *UserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	rows, err := u.db.Query(ctx,
		`SELECT id, full_name, email, password, phone FROM users`)
	if err != nil {
		fmt.Printf("[DEBUG] Query error: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Phone)
		if err != nil {
			fmt.Printf("[DEBUG] Scan error: %v\n", err)
			return nil, err
		}
		fmt.Printf("[DEBUG] Retrieved user: ID=%d, Name=%s, Email=%s\n", user.ID, user.Name, user.Email)
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		fmt.Printf("[DEBUG] Rows iteration error: %v\n", err)
		return nil, err
	}

	fmt.Printf("[DEBUG] Total users retrieved: %d\n", len(users))
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

// func (u *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
// 	var user models.User

// 	err := u.db.QueryRow(ctx,
// 		`SELECT id, full_name, email, password FROM users WHERE email = $1`,
// 		email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &user, nil
// }

func (u *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	err := u.db.QueryRow(ctx,
		`SELECT u.id, u.full_name, u.email, u.password, ur.role 
		 FROM users u
		 LEFT JOIN user_roles ur ON u.id = ur.user_id
		 WHERE u.email = $1`,
		email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepository) Create(ctx context.Context, user *models.User) error {
	tx, err := u.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx,
		`INSERT INTO users (full_name, email, password) 
		 VALUES ($1, $2, $3)
		 RETURNING id`,
		user.Name, user.Email, user.Password).Scan(&user.ID)
	
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO user_roles (user_id, role)
		 VALUES ($1, $2)`,
		user.ID, user.Role)
	
	if err != nil {
		return fmt.Errorf("failed to assign role: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (u *UserRepository) Update(ctx context.Context, user *models.User) error {
	_, err := u.db.Exec(ctx,
		`UPDATE users SET full_name = $1, email = $2, password = $3, phone = $4 
		 WHERE id = $5`,
		user.Name, user.Email, user.Password, user.Phone, user.ID)

	return err
}

func (u *UserRepository) UpdatePassword(ctx context.Context, password string, id int) error {
	_, err := u.db.Exec(ctx,
		`UPDATE users SET password = $1
		 WHERE id = $2`,
		password, id)

	return err
}

func (u *UserRepository) Delete(ctx context.Context, idInt int) error {
	_, err := u.db.Exec(ctx,
		`DELETE FROM users WHERE id = $1`,
		idInt)

	return err
}

func (u *UserRepository) AssignRole(ctx context.Context, userID int, role string) error {
	_, err := u.db.Exec(ctx,
		`UPDATE user_roles SET role = $1, updated_at = CURRENT_TIMESTAMP
		 WHERE user_id = $2`,
		role, userID)

	return err
}