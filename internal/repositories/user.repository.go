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
		username,
		email, 
		mobile_no,
		password_hash,
		full_name,
	) 
	VALUES($1, $2, $3, $4, $5, $6)`
	tag, err := userRepo.pool.Exec(c, userCreate,
		user.Username,
		user.Email,
		user.MobileNo,
		user.PasswordHash,
		user.FullName,
	)
	if err != nil {
		return err
	}
	log.Printf("Insert data into users table. No. of Rows affected: %v", tag.RowsAffected())
	return nil
}
