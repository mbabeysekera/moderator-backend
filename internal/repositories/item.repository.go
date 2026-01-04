package repositories

import (
	"context"
	"errors"
	"log/slog"

	"coolbreez.lk/moderator/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ItemRepository struct {
	pool *pgxpool.Pool
}

var ErrItemNotAffected = errors.New("no rows affected")

func NewItemRepository(dbPool *pgxpool.Pool) *ItemRepository {
	return &ItemRepository{
		pool: dbPool,
	}
}

func (ir *ItemRepository) Create(ctx context.Context, item *models.Item) error {
	const createItem = `INSERT INTO items (
		product_id,
		image_url,
	)
	VALUES($1, $2)
	`
	tag, err := ir.pool.Exec(ctx, createItem, item.ProductID, item.ImageURL)
	if err != nil {
		slog.Error("db insert",
			"repository", "item",
			"err", err,
			"query", createItem,
			"product_id", item.ProductID,
		)
		return err
	}
	if tag.RowsAffected() == 0 {
		slog.Warn("db insert items details",
			"repository", "item",
			"err", ErrRowsNotAffected,
			"query", createItem,
			"user_id", nil,
		)
		return ErrRowsNotAffected
	}
	return nil
}
