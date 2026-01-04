package services

import (
	"context"
	"errors"
	"log/slog"
	"strconv"

	enums "coolbreez.lk/moderator/internal/constants"
	"coolbreez.lk/moderator/internal/dto"
	"coolbreez.lk/moderator/internal/middlewares"
	"coolbreez.lk/moderator/internal/models"
	"coolbreez.lk/moderator/internal/repositories"
	"coolbreez.lk/moderator/internal/utils"
)

type ProductServiceImpl struct {
	productRepo *repositories.ProductRepository
}

var ErrProductItemCreateFailed = errors.New("product items create failed")
var ErrProductItemUpdateFailed = errors.New("product items update failed")

func NewProductService(repo *repositories.ProductRepository) *ProductServiceImpl {
	return &ProductServiceImpl{
		productRepo: repo,
	}
}

func (ps *ProductServiceImpl) CreateProductWithItems(c context.Context,
	productsWithItems *dto.ProductsWithItemsRequest) error {

	addedBy := c.Value(middlewares.AuthorizationContextKey).(*utils.JWTExtractedDetails).UserID

	product := &models.Product{
		Title:       productsWithItems.Title,
		Brand:       productsWithItems.Brand,
		Category:    productsWithItems.Category,
		Sku:         productsWithItems.Sku,
		Description: productsWithItems.Description,
		Price:       productsWithItems.Price,
		AddedBy:     addedBy,
		InStock:     productsWithItems.InStock,
	}
	items := make([]models.Item, 0)
	for _, item := range productsWithItems.Items {
		item := &models.Item{
			ImageURL: item.ImageURL,
		}
		items = append(items, *item)
	}
	err := ps.productRepo.CreateProductWithItems(c, product, items)
	if err != nil {
		slog.Error("products details create",
			"service", "product",
			"err", err,
			"action", "create",
			"added_by", addedBy,
		)
		if errors.Is(err, repositories.ErrRowsNotAffected) {
			return ErrProductItemCreateFailed
		}
		return err
	}
	slog.Info("products details created",
		"service", "product",
		"action", "create",
		"added_by", addedBy,
	)
	return nil
}

func (ps *ProductServiceImpl) GetProductsWithItems(c context.Context,
	count string, page string, category enums.ProductCategory) ([]repositories.ProductWithItems, error) {

	limit, err := strconv.ParseInt(count, 10, 64)
	if err != nil {
		slog.Warn("products details fetch",
			"service", "product",
			"err", err,
			"action", "fetch",
		)
		return nil, ErrProductFetch
	}
	pageNo, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		slog.Warn("products details fetch",
			"service", "product",
			"err", err,
			"action", "fetch",
		)
		return nil, ErrProductFetch
	}
	offset := (pageNo - 1) * limit
	productsWithItems, err := ps.productRepo.GetProductsWithItems(c, limit, offset, category)
	if err != nil {
		slog.Error("products details fetch",
			"service", "product",
			"err", err,
			"action", "fetch",
		)
		return nil, err
	}
	slog.Info("products details fetched",
		"service", "product",
		"action", "fetch",
	)
	return productsWithItems, nil
}

func (ps *ProductServiceImpl) GetProductWithItems(c context.Context,
	id string) (*repositories.ProductWithItems, error) {
	productID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		slog.Warn("products details fetch",
			"service", "product",
			"err", err,
			"action", "fetch",
		)
		return nil, ErrInvalidParams
	}
	productWithItems, err := ps.productRepo.GetProductWithItemsByID(c, productID)
	if productWithItems == nil {
		slog.Info("product details fetch",
			"service", "product",
			"action", "fetch",
			"product_id", productID,
		)
		return nil, ErrInvalidProduct
	}
	if err != nil {
		slog.Error("product details fetch",
			"service", "product",
			"err", err,
			"action", "fetch",
		)
		return nil, err
	}
	slog.Info("product details fetched",
		"service", "product",
		"action", "fetch",
	)
	return productWithItems, nil
}

func (ps *ProductServiceImpl) DeleteProductByID(c context.Context, id string) error {

	productID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return ErrInvalidParams
	}

	deleteRequestedBy := c.Value(middlewares.AuthorizationContextKey).(*utils.JWTExtractedDetails).UserID

	err = ps.productRepo.DeleteProductByID(c, productID, deleteRequestedBy)
	if err != nil {
		slog.Error("products details delete",
			"service", "product",
			"err", err,
			"action", "delete",
		)
		if errors.Is(err, repositories.ErrRowsNotAffected) {
			return ErrProductDelete
		}
		return ErrProductFetch
	}
	slog.Info("products details deleted",
		"service", "product",
		"action", "delete",
	)
	return nil
}

func (ps *ProductServiceImpl) GetProductWithItemsBySku(c context.Context,
	sku string) (*repositories.ProductWithItems, error) {
	productWithItems, err := ps.productRepo.GetProductBySku(c, sku)
	if productWithItems == nil {
		slog.Info("product details fetch",
			"service", "product",
			"action", "fetch",
			"product_sku", sku,
		)
		return nil, ErrInvalidProduct
	}
	if err != nil {
		slog.Error("product details fetch",
			"service", "product",
			"err", err,
			"action", "fetch",
		)
		return nil, err
	}
	slog.Info("product details fetched",
		"service", "product",
		"action", "fetch",
	)
	return productWithItems, nil
}

func (ps *ProductServiceImpl) UpdateProductStock(c context.Context,
	stock int, productID int64) error {

	updatedBy := c.Value(middlewares.AuthorizationContextKey).(*utils.JWTExtractedDetails).UserID

	err := ps.productRepo.UpdateProductStockByID(c, stock, productID)
	if err != nil {
		slog.Error("products details update",
			"service", "product",
			"err", err,
			"action", "update",
			"added_by", updatedBy,
		)
		if errors.Is(err, repositories.ErrRowsNotAffected) {
			return ErrProductItemUpdateFailed
		}
		return err
	}
	slog.Info("products details update",
		"service", "product",
		"action", "update",
		"added_by", updatedBy,
	)
	return nil
}

func (ps *ProductServiceImpl) UpdateProductPrice(c context.Context,
	price float64, productID int64) error {

	updatedBy := c.Value(middlewares.AuthorizationContextKey).(*utils.JWTExtractedDetails).UserID

	err := ps.productRepo.UpdateProductPriceByID(c, price, productID)
	if err != nil {
		slog.Error("products details update",
			"service", "product",
			"err", err,
			"action", "update",
			"added_by", updatedBy,
		)
		if errors.Is(err, repositories.ErrRowsNotAffected) {
			return ErrProductItemUpdateFailed
		}
		return err
	}
	slog.Info("products details update",
		"service", "product",
		"action", "update",
		"added_by", updatedBy,
	)
	return nil
}
