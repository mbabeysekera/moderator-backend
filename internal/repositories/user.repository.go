package repositories

import (
	"context"
	"log"

	"coolbreez.lk/moderator/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(dbPool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: dbPool,
	}
}

func (userRepo *UserRepository) Create(c context.Context, user *models.User) error {
	userCreate := `INSERT INTO users(
		email, 
		mobile_no,
		password_hash,
		full_name,
		role
	) 
	VALUES($1, $2, $3, $4, $5)`
	tag, err := userRepo.pool.Exec(c, userCreate,
		user.Email,
		user.MobileNo,
		user.PasswordHash,
		user.FullName,
		user.Role,
	)
	if err != nil {
		return err
	}
	log.Printf("Insert data into users table. No. of Rows affected: %v", tag.RowsAffected())
	return nil
}

func (userRepo *UserRepository) GetUserByMobileNo(c context.Context,
	mobileNo string) (*models.User, error) {
	getUser := `SELECT 
		id,
		email, 
		mobile_no,
		password_hash,
		full_name,
		role,
		is_active,
		failed_login_attempts,
		last_login_at,
		created_at,
		updated_at
	FROM users WHERE mobile_no = $1
	`
	userRow := userRepo.pool.QueryRow(c, getUser, mobileNo)
	var user models.User
	err := userRow.Scan(&user)
	if err != nil {
		return nil, err
	}
	return &user, err
}
