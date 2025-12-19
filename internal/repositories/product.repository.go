package repositories

import (
	"context"
	"errors"
	"log/slog"

	"coolbreez.lk/moderator/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	pool *pgxpool.Pool
}

type ProductWithItems struct {
	Product models.Product `json:"product"`
	Items   []models.Item  `json:"items"`
}

func NewProductRepository(dbPool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{
		pool: dbPool,
	}
}

func (pr *ProductRepository) CreateProductWithItems(ctx context.Context,
	product *models.Product, items []models.Item) error {
	tx, err := pr.pool.Begin(ctx)
	if err != nil {
		slog.Error("db transaction",
			"repository", "product",
			"err", err,
			"query", "",
			"added_by", product.AddedBy,
		)
		return err
	}
	defer tx.Rollback(ctx)

	const addProduct = `INSERT INTO  
		products(title, brand, category, sku, description, price, added_by)
		VALUES($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	err = tx.QueryRow(ctx, addProduct,
		product.Title,
		product.Brand,
		product.Category,
		product.Sku,
		product.Description,
		product.Price,
		product.AddedBy,
	).Scan(&product.ID)
	if err != nil {
		slog.Error("db insert",
			"repository", "product",
			"err", err,
			"query", addProduct,
			"added_by", product.AddedBy,
		)
		return err
	}
	const addItems = `INSERT INTO
		items(product_id, item_code, in_stock, image_url)
		VALUES($1, $2, $3, $4)
	`
	for _, item := range items {
		tag, err := tx.Exec(ctx, addItems,
			product.ID,
			item.ItemCode,
			item.InStock,
			item.ImageURL,
		)
		if err != nil {
			slog.Error("db insert",
				"repository", "item",
				"err", err,
				"query", addItems,
				"product_id", product.ID,
			)
			return err
		}
		if tag.RowsAffected() == 0 {
			slog.Warn("db insert item details",
				"repository", "item",
				"err", ErrRowsNotAffected,
				"query", addItems,
				"product_id", product.ID,
			)
			return ErrRowsNotAffected
		}
	}
	return tx.Commit(ctx)
}

func (pr *ProductRepository) GetProductsWithItems(ctx context.Context,
	limit int64, offset int64) ([]ProductWithItems, error) {

	const getProducts = ` SELECT id, title, brand, category, sku, description, price, created_at
			FROM products ORDER BY created_at DESC LIMIT $1 OFFSET $2
		`
	rows, err := pr.pool.Query(ctx, getProducts,
		limit,
		offset,
	)
	if err != nil {
		slog.Error("db fetch",
			"repository", "product",
			"err", err,
			"query", getProducts,
		)
		return nil, ErrDBQuery
	}
	defer rows.Close()

	var productsWithItems = make([]ProductWithItems, 0)
	var productsWithItemsMap = make(map[int64]ProductWithItems)
	var productIDs = make([]int64, 0)
	for rows.Next() {
		var productWithItems ProductWithItems
		err = rows.Scan(
			&productWithItems.Product.ID,
			&productWithItems.Product.Title,
			&productWithItems.Product.Brand,
			&productWithItems.Product.Category,
			&productWithItems.Product.Sku,
			&productWithItems.Product.Description,
			&productWithItems.Product.Price,
			&productWithItems.Product.CreatedAt,
		)
		if err != nil {
			slog.Error("db rows",
				"repository", "product",
				"err", err,
			)
			return nil, ErrDBQuery
		}
		productWithItems.Items = make([]models.Item, 0)
		productsWithItemsMap[productWithItems.Product.ID] = productWithItems
		productIDs = append(productIDs, productWithItems.Product.ID)
	}
	err = rows.Err()
	if err != nil {
		slog.Error("db rows",
			"repository", "product",
			"err", err,
		)
		return nil, ErrDBQuery
	}

	if len(productIDs) == 0 {
		return productsWithItems, nil
	}

	const getItemsForProduct = `SELECT id, product_id, item_code, in_stock, image_url, created_at
		FROM items WHERE product_id = ANY($1)
	`
	itemRows, err := pr.pool.Query(ctx, getItemsForProduct, productIDs)
	if err != nil {
		slog.Error("db fetch",
			"repository", "items",
			"err", err,
			"query", getItemsForProduct,
		)
		return nil, ErrDBQuery
	}
	defer itemRows.Close()

	for itemRows.Next() {
		var item models.Item
		err = itemRows.Scan(
			&item.ID,
			&item.ProductID,
			&item.ItemCode,
			&item.InStock,
			&item.ImageURL,
			&item.CreatedAt,
		)
		if err != nil {
			slog.Error("db rows",
				"repository", "items",
				"err", err,
			)
			return nil, ErrDBQuery
		}
		product := productsWithItemsMap[item.ProductID]
		product.Items = append(productsWithItemsMap[item.ProductID].Items, item)
		productsWithItemsMap[item.ProductID] = product
	}
	err = itemRows.Err()
	if err != nil {
		slog.Error("db rows",
			"repository", "items",
			"err", err,
		)
		return nil, ErrDBQuery
	}
	for _, prod := range productsWithItemsMap {
		productsWithItems = append(productsWithItems, prod)
	}
	return productsWithItems, nil
}

func (pr *ProductRepository) GetProductWithItemsByID(ctx context.Context,
	productID int64) (*ProductWithItems, error) {
	const getProductByID = ` SELECT id, title, brand, category, sku, description, price, created_at
			FROM products WHERE id = $1
		`
	productRow := pr.pool.QueryRow(ctx, getProductByID, productID)
	var productWithItems ProductWithItems
	err := productRow.Scan(
		&productWithItems.Product.ID,
		&productWithItems.Product.Title,
		&productWithItems.Product.Brand,
		&productWithItems.Product.Category,
		&productWithItems.Product.Sku,
		&productWithItems.Product.Description,
		&productWithItems.Product.Price,
		&productWithItems.Product.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		slog.Error("db rows",
			"repository", "product",
			"err", err,
		)
		return nil, ErrDBQuery
	}
	productWithItems.Items = make([]models.Item, 0)
	const getItemsForProduct = `SELECT id, product_id, item_code, in_stock, image_url, created_at
		FROM items WHERE product_id = $1
	`
	itemsRows, err := pr.pool.Query(ctx, getItemsForProduct, productID)
	if err != nil {
		slog.Error("db fetch",
			"repository", "items",
			"err", err,
		)
		return nil, ErrDBQuery
	}
	defer itemsRows.Close()

	for itemsRows.Next() {
		var item models.Item
		itemsRows.Scan(
			&item.ID,
			&item.ProductID,
			&item.ItemCode,
			&item.InStock,
			&item.ImageURL,
			&item.CreatedAt,
		)
		productWithItems.Items = append(productWithItems.Items, item)
	}
	err = itemsRows.Err()
	if err != nil {
		slog.Error("db rows",
			"repository", "items",
			"err", err,
		)
		return nil, ErrDBQuery
	}
	return &productWithItems, nil
}

func (pr *ProductRepository) DeleteProductByID(ctx context.Context,
	productID int64, deleteRequestedBy int64) error {

	const deleteProduct = `DELETE FROM products WHERE id = $1 AND added_by = $2`

	tag, err := pr.pool.Exec(ctx, deleteProduct, productID, deleteRequestedBy)
	if err != nil {
		slog.Error("db delete",
			"repository", "products,items",
			"err", err,
			"query", deleteProduct,
			"product_id", productID,
		)
		return ErrDBQuery
	}
	if tag.RowsAffected() == 0 {
		slog.Warn("db delete products,items details",
			"repository", "products,items",
			"err", ErrRowsNotAffected,
			"query", deleteProduct,
			"product_id", productID,
		)
		return ErrRowsNotAffected
	}
	return nil
}
