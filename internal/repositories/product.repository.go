package repositories

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	enums "coolbreez.lk/moderator/internal/constants"
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
			"app_id", product.AppID,
		)
		return err
	}
	defer tx.Rollback(ctx)

	const addProduct = `INSERT INTO  
		products(title, brand, category, sku, description, price, added_by, in_stock, app_id)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
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
		product.InStock,
		product.AppID,
	).Scan(&product.ID)
	if err != nil {
		slog.Error("db insert",
			"repository", "product",
			"err", err,
			"query", addProduct,
			"added_by", product.AddedBy,
			"app_id", product.AppID,
		)
		return err
	}
	const addItems = `INSERT INTO
		items(product_id, image_url, app_id)
		VALUES($1, $2, $3)
	`
	for _, item := range items {
		tag, err := tx.Exec(ctx, addItems,
			product.ID,
			item.ImageURL,
			product.AppID,
		)
		if err != nil {
			slog.Error("db insert",
				"repository", "item",
				"err", err,
				"query", addItems,
				"product_id", product.ID,
				"app_id", product.AppID,
			)
			return err
		}
		if tag.RowsAffected() == 0 {
			slog.Warn("db insert item details",
				"repository", "item",
				"err", ErrRowsNotAffected,
				"query", addItems,
				"product_id", product.ID,
				"app_id", product.AppID,
			)
			return ErrRowsNotAffected
		}
	}
	return tx.Commit(ctx)
}

func (pr *ProductRepository) GetProductsWithItems(ctx context.Context,
	limit int64, offset int64, category enums.ProductCategory, appID int64) ([]ProductWithItems, error) {

	getProducts := `SELECT id, title, brand, category, sku, description, price, created_at, in_stock, app_id
			FROM products WHERE app_id = $1
		`
	args := []any{appID}
	argPosition := 2

	if category != enums.ProductCategory("ALL") {
		getProducts += fmt.Sprintf(" AND category = $%d", argPosition)
		args = append(args, string(category))
		argPosition++
	}

	getProducts += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argPosition, argPosition+1)
	args = append(args, limit, offset)

	rows, err := pr.pool.Query(ctx, getProducts,
		args...,
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
			&productWithItems.Product.InStock,
			&productWithItems.Product.AppID,
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

	const getItemsForProduct = `SELECT id, product_id, image_url, created_at, app_id
		FROM items WHERE product_id = ANY($1) AND app_id = $2
	`
	itemRows, err := pr.pool.Query(ctx, getItemsForProduct, productIDs, appID)
	if err != nil {
		slog.Error("db fetch",
			"repository", "items",
			"err", err,
			"query", getItemsForProduct,
			"product_ids", productIDs,
			"app_id", appID,
		)
		return nil, ErrDBQuery
	}
	defer itemRows.Close()

	for itemRows.Next() {
		var item models.Item
		err = itemRows.Scan(
			&item.ID,
			&item.ProductID,
			&item.ImageURL,
			&item.CreatedAt,
			&item.AppID,
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
	productID int64, appID int64) (*ProductWithItems, error) {
	const getProductByID = ` SELECT id, title, brand, category, sku, description, price, created_at,
			in_stock, app_id 
			FROM products WHERE id = $1 AND app_id = $2
		`
	productRow := pr.pool.QueryRow(ctx, getProductByID, productID, appID)
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
		&productWithItems.Product.InStock,
		&productWithItems.Product.AppID,
	)
	if err != nil {
		slog.Error("db rows",
			"repository", "product",
			"err", err,
		)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, ErrDBQuery
	}
	productWithItems.Items = make([]models.Item, 0)
	const getItemsForProduct = `SELECT id, product_id, image_url, created_at, app_id
		FROM items WHERE product_id = $1 AND app_id = $2
	`
	itemsRows, err := pr.pool.Query(ctx, getItemsForProduct, productID, appID)
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
			&item.ImageURL,
			&item.CreatedAt,
			&item.AppID,
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
	productID int64, appID int64, deleteRequestedBy int64) error {

	const deleteProduct = `DELETE FROM products WHERE id = $1 AND added_by = $2 AND app_id = $3`

	tag, err := pr.pool.Exec(ctx, deleteProduct, productID, deleteRequestedBy, appID)
	if err != nil {
		slog.Error("db delete",
			"repository", "products,items",
			"err", err,
			"query", deleteProduct,
			"product_id", productID,
			"app_id", appID,
			"delete_requested_by", deleteRequestedBy,
		)
		return ErrDBQuery
	}
	if tag.RowsAffected() == 0 {
		slog.Warn("db delete products,items details",
			"repository", "products,items",
			"err", ErrRowsNotAffected,
			"query", deleteProduct,
			"product_id", productID,
			"app_id", appID,
			"delete_requested_by", deleteRequestedBy,
		)
		return ErrRowsNotAffected
	}
	return nil
}

func (pr *ProductRepository) GetProductBySku(ctx context.Context,
	sku string, appID int64) (*ProductWithItems, error) {

	const getProductBySku = `SELECT id, title, brand, category, sku, description, price, created_at,
			in_stock, app_id
			FROM products WHERE sku = $1 AND app_id = $2
		`
	productRow := pr.pool.QueryRow(ctx, getProductBySku, sku, appID)
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
		&productWithItems.Product.InStock,
		&productWithItems.Product.AppID,
	)
	if err != nil {
		slog.Error("db rows",
			"repository", "product",
			"err", err,
		)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, ErrDBQuery
	}
	productWithItems.Items = make([]models.Item, 0)
	const getItemsForProduct = `SELECT id, product_id, image_url, created_at, app_id
		FROM items WHERE product_id = $1 AND app_id = $2
	`
	itemsRows, err := pr.pool.Query(ctx, getItemsForProduct, productWithItems.Product.ID, appID)
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
			&item.ImageURL,
			&item.CreatedAt,
			&item.AppID,
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

func (pr *ProductRepository) GetTotalProductsCount(ctx context.Context,
	category enums.ProductCategory, appID int64) (int64, error) {

	countProducts := `SELECT count(*) FROM products WHERE app_id = $1`

	args := []any{appID}
	argPos := 2

	if category != enums.ProductCategory("ALL") {
		countProducts += fmt.Sprintf(" AND category = $%d", argPos)
		args = append(args, category)
	}

	rowCount := pr.pool.QueryRow(ctx, countProducts, args...)
	var count int64
	err := rowCount.Scan(&count)
	if err != nil {
		slog.Error("db fetch",
			"repository", "product",
			"err", err,
			"query", countProducts,
		)
		return 0, ErrDBQuery
	}

	return count, nil
}

func (pr *ProductRepository) UpdateProductStockByID(ctx context.Context,
	stock int, productID int64, appID int64) error {
	const updateProduct = `UPDATE products SET in_stock = $1 WHERE id = $2 AND app_id = $3`
	tag, err := pr.pool.Exec(ctx, updateProduct, stock, productID, appID)
	if err != nil {
		slog.Error("db update",
			"repository", "product",
			"err", err,
			"query", updateProduct,
			"product_id", productID,
		)
		return err
	}
	if tag.RowsAffected() == 0 {
		slog.Warn("db insert item details",
			"repository", "product",
			"err", ErrRowsNotAffected,
			"query", updateProduct,
			"product_id", productID,
		)
		return ErrRowsNotAffected
	}
	return nil
}

func (pr *ProductRepository) UpdateProductPriceByID(ctx context.Context,
	price float64, productID int64, appID int64) error {
	const updateProduct = `UPDATE products SET price = $1 WHERE id = $2 AND app_id = $3`
	tag, err := pr.pool.Exec(ctx, updateProduct, price, productID, appID)
	if err != nil {
		slog.Error("db update",
			"repository", "product",
			"err", err,
			"query", updateProduct,
			"product_id", productID,
		)
		return err
	}
	if tag.RowsAffected() == 0 {
		slog.Warn("db insert item details",
			"repository", "product",
			"err", ErrRowsNotAffected,
			"query", updateProduct,
			"product_id", productID,
		)
		return ErrRowsNotAffected
	}
	return nil
}
