package repositories

import (
	"context"
	"errors"
	"log/slog"

	"coolbreez.lk/moderator/internal/models"
	"github.com/jackc/pgx/v5"
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

func (userRepo *UserRepository) Create(ctx context.Context, user *models.User) error {
	const userCreate = `INSERT INTO users(
		email, 
		mobile_no,
		password_hash,
		full_name,
		role
	) 
	VALUES($1, $2, $3, $4, $5)`
	tag, err := userRepo.pool.Exec(ctx, userCreate,
		user.Email,
		user.MobileNo,
		user.PasswordHash,
		user.FullName,
		user.Role,
	)
	if err != nil {
		slog.Error("db insert",
			"repository", "user",
			"err", err,
			"query", userCreate,
			"mobile_no", user.MobileNo,
		)
		return err
	}
	if tag.RowsAffected() == 0 {
		slog.Warn("db insert user details",
			"repository", "user",
			"err", ErrRowsNotAffected,
			"query", userCreate,
			"user_id", nil,
		)
		return ErrRowsNotAffected
	}
	return nil
}

func (userRepo *UserRepository) GetUserByMobileNo(ctx context.Context,
	mobileNo string) (*models.User, error) {
	const getUser = `SELECT 
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
	userRow := userRepo.pool.QueryRow(ctx, getUser, mobileNo)
	var user models.User
	err := userRow.Scan(
		&user.ID,
		&user.Email,
		&user.MobileNo,
		&user.PasswordHash,
		&user.FullName,
		&user.Role,
		&user.IsActive,
		&user.FailedLoginAttempts,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		slog.Error("db select user by modile",
			"repository", "user",
			"err", err,
			"query", getUser,
			"mobile_no", mobileNo,
		)
		return nil, err
	}
	return &user, nil
}

func (userRepo *UserRepository) IncrementUserLoginFailuresByID(ctx context.Context,
	userID int64) error {
	const updateFailedLogin = `UPDATE
	users SET
		failed_login_attempts = failed_login_attempts + 1,
		updated_at = NOW()
	WHERE id = $1	
	`
	tag, err := userRepo.pool.Exec(ctx, updateFailedLogin, userID)
	if err != nil {
		slog.Error("db update user for login failures",
			"repository", "user",
			"err", err,
			"query", updateFailedLogin,
			"user_id", userID,
		)
		return err
	}
	if tag.RowsAffected() == 0 {
		slog.Warn("db update user details",
			"repository", "user",
			"err", ErrRowsNotAffected,
			"query", updateFailedLogin,
			"user_id", userID,
		)
		return ErrRowsNotAffected
	}
	return nil
}

func (userRepo *UserRepository) UpdateSuccessfulLoginByID(ctx context.Context,
	userID int64) error {
	const updateLogin = `UPDATE
		users SET
			failed_login_attempts = 0,
			last_login_at = NOW(),
			updated_at = NOW()
		WHERE id = $1
		`
	tag, err := userRepo.pool.Exec(ctx, updateLogin, userID)
	if err != nil {
		slog.Error("db update user for login success",
			"repository", "user",
			"err", err,
			"query", updateLogin,
			"user_id", userID,
		)
		return err
	}
	if tag.RowsAffected() == 0 {
		slog.Warn("db update user details",
			"repository", "user",
			"err", ErrRowsNotAffected,
			"query", updateLogin,
			"user_id", userID,
		)
		return ErrRowsNotAffected
	}
	return nil
}

func (userRepo *UserRepository) UpdateUserByID(ctx context.Context, user *models.User) error {
	const userUpdate = `UPDATE
	users SET
		email = $1,
		mobile_no = $2,
		full_name = $3,
		updated_at = NOW()
	WHERE id = $4
	`
	tag, err := userRepo.pool.Exec(ctx, userUpdate, user.Email, user.MobileNo, user.FullName, user.ID)
	if err != nil {
		slog.Error("db update user details",
			"repository", "user",
			"err", err,
			"query", userUpdate,
			"user_id", user.ID,
		)
		return err
	}
	if tag.RowsAffected() == 0 {
		slog.Warn("db update user details",
			"repository", "user",
			"err", ErrRowsNotAffected,
			"query", userUpdate,
			"user_id", user.ID,
		)
		return ErrRowsNotAffected
	}
	return nil
}

func (userRepo *UserRepository) GetUserByID(ctx context.Context,
	userID int64) (*models.User, error) {
	const getUserByID = `SELECT id, full_name, role FROM users WHERE id = $1`
	userRow := userRepo.pool.QueryRow(ctx, getUserByID, userID)
	var user models.User
	err := userRow.Scan(
		&user.ID,
		&user.FullName,
		&user.Role,
	)
	if err != nil {
		slog.Error("db update user details",
			"repository", "user",
			"err", err,
			"query", getUserByID,
			"user_id", userID,
		)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, ErrDBQuery
	}
	return &user, nil
}
