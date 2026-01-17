package controllers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	enums "coolbreez.lk/moderator/internal/constants"
	"coolbreez.lk/moderator/internal/dto"
	apperrors "coolbreez.lk/moderator/internal/errors"
	"coolbreez.lk/moderator/internal/repositories"
	"coolbreez.lk/moderator/internal/services"
	"github.com/gin-gonic/gin"
)

type ProductService interface {
	CreateProductWithItems(c context.Context,
		productsWithItems *dto.ProductsWithItemsRequest) error
	GetProductsWithItems(c context.Context,
		count string, page string, appID string,
		category enums.ProductCategory) ([]repositories.ProductWithItems, error)
	GetProductWithItems(c context.Context,
		id string, appID string) (*repositories.ProductWithItems, error)
	DeleteProductByID(c context.Context, id string, appID string) error
	GetProductWithItemsBySku(c context.Context,
		sku string, appID string) (*repositories.ProductWithItems, error)
	UpdateProductStock(c context.Context,
		stock int, productID int64, appID string) error
	UpdateProductPrice(c context.Context,
		price float64, productID int64, appID string) error
}

type ProductController struct {
	service ProductService
}

func NewProductController(productService ProductService) *ProductController {
	return &ProductController{
		service: productService,
	}
}

func (pc *ProductController) CreateProductWithItems(c *gin.Context) {
	var productWithItems dto.ProductsWithItemsRequest
	err := c.ShouldBindJSON(&productWithItems)
	if err != nil {
		slog.Error("productWithItems parameter validation",
			"err", err,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"ip", c.ClientIP(),
		)
		c.JSON(http.StatusBadRequest,
			apperrors.AppStdErrorHandler(
				"parameter validation failed",
				"ps_0000",
			),
		)
		return
	}
	err = pc.service.CreateProductWithItems(c.Request.Context(), &productWithItems)
	if err != nil {
		if errors.Is(err, services.ErrProductItemCreateFailed) {
			c.JSON(http.StatusConflict,
				apperrors.AppStdErrorHandler(
					services.ErrProductItemCreateFailed.Error(),
					"ps_0001",
				),
			)
			return
		}
		c.JSON(http.StatusInternalServerError,
			apperrors.AppStdErrorHandler(
				"Internal server error",
				"ps_0001",
			),
		)
		return
	}
	c.JSON(http.StatusOK, &dto.SuccessStdResponse{
		Status:  enums.RequestSuccess,
		Message: "products with items created",
		Time:    time.Now().UTC(),
	})
}

func (pc *ProductController) GetProductsWithItems(c *gin.Context) {
	count := c.DefaultQuery("count", "10")
	page := c.DefaultQuery("page", "1")
	category := c.DefaultQuery("category", "ALL")
	appID, _ := c.Get("app_id")
	productsWithItems, err := pc.service.GetProductsWithItems(c.Request.Context(), count, page,
		appID.(string), enums.ProductCategory(category))
	if err != nil {
		if errors.Is(err, services.ErrProductFetch) {
			c.JSON(http.StatusInternalServerError,
				apperrors.AppStdErrorHandler(
					services.ErrProductFetch.Error(),
					"ps_0000",
				),
			)
			return
		}
		c.JSON(http.StatusInternalServerError,
			apperrors.AppStdErrorHandler(
				"Internal server error",
				"ps_0001",
			),
		)
		return
	}
	c.JSON(http.StatusOK, &dto.ProductsWithItemsResponse{
		All:   productsWithItems,
		Count: len(productsWithItems),
	})
}

func (pc *ProductController) GetProductWithItems(c *gin.Context) {
	productID := c.Param("product_id")
	appID, _ := c.Get("app_id")
	productsWithItems, err := pc.service.GetProductWithItems(c.Request.Context(), productID, appID.(string))
	if err != nil {
		if errors.Is(err, services.ErrInvalidProduct) {
			c.JSON(http.StatusBadRequest,
				apperrors.AppStdErrorHandler(
					services.ErrInvalidProduct.Error(),
					"ps_0000",
				),
			)
			return
		}
		if errors.Is(err, services.ErrInvalidParams) {
			c.JSON(http.StatusBadRequest,
				apperrors.AppStdErrorHandler(
					services.ErrInvalidParams.Error(),
					"ps_0001",
				),
			)
			return
		}
		c.JSON(http.StatusInternalServerError,
			apperrors.AppStdErrorHandler(
				"Internal server error",
				"ps_0002",
			),
		)
		return
	}
	c.JSON(http.StatusOK, &dto.ProductWithItemsResponse{
		Product: productsWithItems.Product,
		Items:   productsWithItems.Items,
	})
}

func (pc *ProductController) DeleteProductByID(c *gin.Context) {
	productID := c.Param("product_id")
	appID, _ := c.Get("app_id")
	err := pc.service.DeleteProductByID(c.Request.Context(), productID, appID.(string))
	if err != nil {
		if errors.Is(err, services.ErrInvalidParams) {
			c.JSON(http.StatusBadRequest,
				apperrors.AppStdErrorHandler(
					services.ErrInvalidParams.Error(),
					"ps_0000",
				),
			)
			return
		}
		if errors.Is(err, services.ErrProductDelete) {
			c.JSON(http.StatusUnauthorized,
				apperrors.AppStdErrorHandler(
					services.ErrProductDelete.Error(),
					"ps_0001",
				),
			)
			return
		}
		c.JSON(http.StatusInternalServerError,
			apperrors.AppStdErrorHandler(
				"Internal server error",
				"ps_0002",
			),
		)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (pc *ProductController) GetProductWithItemsBySku(c *gin.Context) {
	sku := c.Param("sku")
	appID, _ := c.Get("app_id")
	productsWithItems, err := pc.service.GetProductWithItemsBySku(c.Request.Context(), sku, appID.(string))
	if err != nil {
		if errors.Is(err, services.ErrInvalidProduct) {
			c.JSON(http.StatusNotFound,
				apperrors.AppStdErrorHandler(
					services.ErrInvalidProduct.Error(),
					"ps_0000",
				),
			)
			return
		}
		c.JSON(http.StatusInternalServerError,
			apperrors.AppStdErrorHandler(
				"Internal server error",
				"ps_0001",
			),
		)
		return
	}
	c.JSON(http.StatusOK, &dto.ProductWithItemsResponse{
		Product: productsWithItems.Product,
		Items:   productsWithItems.Items,
	})
}

func (pc *ProductController) UpdateProduct(c *gin.Context) {
	var productDetailsToBeUpdated dto.ProductDetailsUpdateRequest
	err := c.ShouldBindJSON(&productDetailsToBeUpdated)
	appID, _ := c.Get("app_id")
	if err != nil {
		slog.Error("productDetailsToBeUpdated parameter validation",
			"err", err,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"ip", c.ClientIP(),
		)
		c.JSON(http.StatusBadRequest,
			apperrors.AppStdErrorHandler(
				"parameter validation failed",
				"ps_0000",
			),
		)
		return
	}
	if productDetailsToBeUpdated.Price != nil {
		errPrice := pc.service.UpdateProductPrice(c.Request.Context(),
			*productDetailsToBeUpdated.Price, productDetailsToBeUpdated.ID, appID.(string))
		if errPrice != nil {
			slog.Error("price update",
				"err", err,
				"method", c.Request.Method,
				"path", c.Request.URL.Path,
				"ip", c.ClientIP(),
			)
			if errors.Is(errPrice, services.ErrProductItemUpdateFailed) {
				c.JSON(http.StatusNotFound,
					apperrors.AppStdErrorHandler(
						services.ErrProductItemUpdateFailed.Error(),
						"ps_0000",
					),
				)
				return
			}
			c.JSON(http.StatusInternalServerError,
				apperrors.AppStdErrorHandler(
					"Internal server error",
					"ps_0001",
				),
			)
			return
		}
	}
	if productDetailsToBeUpdated.InStock != nil {
		errStock := pc.service.UpdateProductStock(c.Request.Context(),
			*productDetailsToBeUpdated.InStock, productDetailsToBeUpdated.ID, appID.(string))
		if errStock != nil {
			slog.Error("price update",
				"err", err,
				"method", c.Request.Method,
				"path", c.Request.URL.Path,
				"ip", c.ClientIP(),
			)
			if errors.Is(errStock, services.ErrProductItemUpdateFailed) {
				c.JSON(http.StatusNotFound,
					apperrors.AppStdErrorHandler(
						services.ErrProductItemUpdateFailed.Error(),
						"ps_0002",
					),
				)
				return
			}
			c.JSON(http.StatusInternalServerError,
				apperrors.AppStdErrorHandler(
					"Internal server error",
					"ps_0003",
				),
			)
			return
		}
	}
	if productDetailsToBeUpdated.Price == nil &&
		productDetailsToBeUpdated.InStock == nil {
		c.JSON(http.StatusBadRequest,
			apperrors.AppStdErrorHandler(
				"price or stock must be set",
				"ps_0004",
			),
		)
		return
	}
	c.JSON(http.StatusOK, &dto.SuccessStdResponse{
		Status:  enums.RequestSuccess,
		Message: "product updated",
		Time:    time.Now().UTC(),
	})
}
