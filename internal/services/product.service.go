package services

import (
	"context"
	"errors"
	"log/slog"
	"strconv"

	"coolbreez.lk/moderator/internal/dto"
	"coolbreez.lk/moderator/internal/middlewares"
	"coolbreez.lk/moderator/internal/models"
	"coolbreez.lk/moderator/internal/repositories"
	"coolbreez.lk/moderator/internal/utils"
)

type ProductServiceImpl struct {
	productRepo *repositories.ProductRepository
}

var ErrProductItemCreateFailed = errors.New("product items failed")

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
		Sku:         productsWithItems.Sku,
		Description: productsWithItems.Description,
		AddedBy:     addedBy,
	}
	items := make([]models.Item, 0)
	for _, item := range productsWithItems.Items {
		item := &models.Item{
			ItemCode: item.ItemCode,
			InStock:  item.InStock,
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
	count string, page string) ([]repositories.ProductsWithItems, error) {

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
	productsWithItems, err := ps.productRepo.GetProductsWithItems(c, limit, offset)
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

func (ps *ProductServiceImpl) DeleteProductByID(c context.Context, productID int64) error {

	deleteRequestedBy := c.Value(middlewares.AuthorizationContextKey).(*utils.JWTExtractedDetails).UserID

	err := ps.productRepo.DeleteProductByID(c, productID, deleteRequestedBy)
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
