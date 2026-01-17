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
	userRole := c.Value(middlewares.AuthorizationContextKey).(*utils.JWTExtractedDetails).UserRole
	appID := c.Value(middlewares.AuthorizationContextKey).(*utils.JWTExtractedDetails).AppID

	if enums.UserRole(userRole) != enums.RoleAdmin && appID != productsWithItems.AppID {
		return ErrProductItemCreateFailed
	}

	product := &models.Product{
		Title:       productsWithItems.Title,
		Brand:       productsWithItems.Brand,
		Category:    productsWithItems.Category,
		Sku:         productsWithItems.Sku,
		Description: productsWithItems.Description,
		Price:       productsWithItems.Price,
		AddedBy:     addedBy,
		InStock:     productsWithItems.InStock,
		AppID:       productsWithItems.AppID,
	}
	items := make([]models.Item, 0)
	for _, item := range productsWithItems.Items {
		item := &models.Item{
			ImageURL: item.ImageURL,
			AppID:    productsWithItems.AppID,
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
	count string, page string, appID string,
	category enums.ProductCategory) ([]repositories.ProductWithItems, error) {

	limit, err := strconv.ParseInt(count, 10, 64)
	if err != nil {
		slog.Warn("products details fetch",
			"service", "product",
			"err", err,
			"action", "validate",
		)
		return nil, ErrInvalidParams
	}
	pageNo, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		slog.Warn("products details fetch",
			"service", "product",
			"err", err,
			"action", "validate",
		)
		return nil, ErrInvalidParams
	}
	offset := (pageNo - 1) * limit
	appIDInt, err := strconv.ParseInt(appID, 10, 64)
	if err != nil {
		slog.Warn("products details fetch",
			"service", "product",
			"err", err,
			"action", "validate",
		)
		return nil, ErrInvalidParams
	}
	productsWithItems, err := ps.productRepo.GetProductsWithItems(c, limit, offset, category, appIDInt)
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
	id string, appID string) (*repositories.ProductWithItems, error) {
	productID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		slog.Warn("products details fetch",
			"service", "product",
			"err", err,
			"action", "validate",
		)
		return nil, ErrInvalidParams
	}
	appIDInt, err := strconv.ParseInt(appID, 10, 64)
	if err != nil {
		slog.Warn("products details fetch",
			"service", "product",
			"err", err,
			"action", "validate",
		)
		return nil, ErrInvalidParams
	}
	productWithItems, err := ps.productRepo.GetProductWithItemsByID(c, productID, appIDInt)
	if productWithItems == nil {
		slog.Info("product details fetch",
			"service", "product",
			"action", "fetch",
			"product_id", productID,
			"app_id", appIDInt,
		)
		return nil, ErrInvalidProduct
	}
	if err != nil {
		slog.Error("product details fetch",
			"service", "product",
			"err", err,
			"action", "fetch",
			"product_id", productID,
			"app_id", appIDInt,
		)
		return nil, err
	}
	slog.Info("product details fetched",
		"service", "product",
		"action", "fetch",
		"product_id", productID,
		"app_id", appIDInt,
	)
	return productWithItems, nil
}

func (ps *ProductServiceImpl) DeleteProductByID(c context.Context, id string, appID string) error {

	productID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		slog.Error("products details delete",
			"service", "product",
			"err", err,
			"action", "validate",
			"product_id", id,
			"app_id", appID,
		)
		return ErrInvalidParams
	}
	appIDInt, err := strconv.ParseInt(appID, 10, 64)
	if err != nil {
		slog.Error("products details delete",
			"service", "product",
			"err", err,
			"action", "validate",
			"product_id", id,
			"app_id", appID,
		)
		return ErrInvalidParams
	}

	deleteRequestedBy := c.Value(middlewares.AuthorizationContextKey).(*utils.JWTExtractedDetails).UserID
	userRole := c.Value(middlewares.AuthorizationContextKey).(*utils.JWTExtractedDetails).UserRole
	appIDfromToken := c.Value(middlewares.AuthorizationContextKey).(*utils.JWTExtractedDetails).AppID

	if enums.UserRole(userRole) != enums.RoleAdmin && appIDfromToken != appIDInt {
		return ErrProductDelete
	}

	err = ps.productRepo.DeleteProductByID(c, productID, appIDInt, deleteRequestedBy)
	if err != nil {
		slog.Error("products details delete",
			"service", "product",
			"err", err,
			"action", "delete",
			"product_id", id,
			"app_id", appID,
			"delete_requested_by", deleteRequestedBy,
		)
		if errors.Is(err, repositories.ErrRowsNotAffected) {
			return ErrProductDelete
		}
		return ErrProductFetch
	}
	slog.Info("products details deleted",
		"service", "product",
		"action", "delete",
		"product_id", id,
		"app_id", appID,
		"delete_requested_by", deleteRequestedBy,
	)
	return nil
}

func (ps *ProductServiceImpl) GetProductWithItemsBySku(c context.Context,
	sku string, appID string) (*repositories.ProductWithItems, error) {
	appIDInt, err := strconv.ParseInt(appID, 10, 64)
	if err != nil {
		slog.Error("products details delete",
			"service", "product",
			"err", err,
			"action", "validate",
			"product_sku", sku,
			"app_id", appID,
		)
		return nil, ErrInvalidParams
	}
	productWithItems, err := ps.productRepo.GetProductBySku(c, sku, appIDInt)
	if productWithItems == nil {
		slog.Info("product details fetch",
			"service", "product",
			"action", "fetch",
			"product_sku", sku,
			"app_id", appID,
		)
		return nil, ErrInvalidProduct
	}
	if err != nil {
		slog.Error("product details fetch",
			"service", "product",
			"err", err,
			"action", "fetch",
			"product_sku", sku,
			"app_id", appID,
		)
		return nil, err
	}
	slog.Info("product details fetched",
		"service", "product",
		"action", "fetch",
		"product_sku", sku,
		"app_id", appID,
	)
	return productWithItems, nil
}

func (ps *ProductServiceImpl) UpdateProductStock(c context.Context,
	stock int, productID int64, appID string) error {

	updatedBy := c.Value(middlewares.AuthorizationContextKey).(*utils.JWTExtractedDetails).UserID
	userRole := c.Value(middlewares.AuthorizationContextKey).(*utils.JWTExtractedDetails).UserRole
	appIDfromToken := c.Value(middlewares.AuthorizationContextKey).(*utils.JWTExtractedDetails).AppID

	appIDInt, err := strconv.ParseInt(appID, 10, 64)
	if err != nil {
		slog.Error("products details update",
			"service", "product",
			"err", err,
			"action", "validate",
			"product_id", productID,
			"app_id", appID,
		)
		return ErrInvalidParams
	}

	if enums.UserRole(userRole) != enums.RoleAdmin && appIDfromToken != appIDInt {
		return ErrProductItemCreateFailed
	}

	err = ps.productRepo.UpdateProductStockByID(c, stock, productID, appIDInt)
	if err != nil {
		slog.Error("products details update",
			"service", "product",
			"err", err,
			"action", "update",
			"added_by", updatedBy,
			"product_id", productID,
			"app_id", appID,
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
		"product_id", productID,
		"app_id", appID,
	)
	return nil
}

func (ps *ProductServiceImpl) UpdateProductPrice(c context.Context,
	price float64, productID int64, appID string) error {
	updatedBy := c.Value(middlewares.AuthorizationContextKey).(*utils.JWTExtractedDetails).UserID
	userRole := c.Value(middlewares.AuthorizationContextKey).(*utils.JWTExtractedDetails).UserRole
	appIDfromToken := c.Value(middlewares.AuthorizationContextKey).(*utils.JWTExtractedDetails).AppID

	appIDInt, err := strconv.ParseInt(appID, 10, 64)
	if err != nil {
		slog.Error("products details update",
			"service", "product",
			"err", err,
			"action", "validate",
			"product_id", productID,
			"app_id", appID,
		)
		return ErrInvalidParams
	}

	if enums.UserRole(userRole) != enums.RoleAdmin && appIDfromToken != appIDInt {
		return ErrProductItemCreateFailed
	}

	err = ps.productRepo.UpdateProductPriceByID(c, price, productID, appIDInt)
	if err != nil {
		slog.Error("products details update",
			"service", "product",
			"err", err,
			"action", "update",
			"added_by", updatedBy,
			"product_id", productID,
			"app_id", appID,
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
		"product_id", productID,
		"app_id", appID,
	)
	return nil
}
