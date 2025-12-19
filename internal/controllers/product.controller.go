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
		count string, page string) ([]repositories.ProductWithItems, error)
	GetProductWithItems(c context.Context,
		id string) (*repositories.ProductWithItems, error)
	DeleteProductByID(c context.Context, id string) error
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
				"us_0000",
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
					"us_0001",
				),
			)
			return
		}
		c.JSON(http.StatusInternalServerError,
			apperrors.AppStdErrorHandler(
				"Internal server error",
				"us_0001",
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
	productsWithItems, err := pc.service.GetProductsWithItems(c.Request.Context(), count, page)
	if err != nil {
		if errors.Is(err, services.ErrProductFetch) {
			c.JSON(http.StatusInternalServerError,
				apperrors.AppStdErrorHandler(
					services.ErrProductFetch.Error(),
					"us_0000",
				),
			)
			return
		}
		c.JSON(http.StatusInternalServerError,
			apperrors.AppStdErrorHandler(
				"Internal server error",
				"us_0001",
			),
		)
		return
	}
	c.JSON(http.StatusOK, &dto.ProductsWithItemsResponse{
		All: productsWithItems,
	})
}

func (pc *ProductController) GetProductWithItems(c *gin.Context) {
	productID := c.Param("product_id")
	productsWithItems, err := pc.service.GetProductWithItems(c.Request.Context(), productID)
	if err != nil {
		if errors.Is(err, services.ErrInvalidProduct) {
			c.JSON(http.StatusBadRequest,
				apperrors.AppStdErrorHandler(
					services.ErrInvalidParams.Error(),
					"us_0000",
				),
			)
			return
		}
		if errors.Is(err, services.ErrInvalidParams) {
			c.JSON(http.StatusBadRequest,
				apperrors.AppStdErrorHandler(
					services.ErrInvalidParams.Error(),
					"us_0001",
				),
			)
			return
		}
		if errors.Is(err, services.ErrProductFetch) {
			c.JSON(http.StatusInternalServerError,
				apperrors.AppStdErrorHandler(
					services.ErrProductFetch.Error(),
					"us_0002",
				),
			)
			return
		}
		c.JSON(http.StatusInternalServerError,
			apperrors.AppStdErrorHandler(
				"Internal server error",
				"us_0003",
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

	err := pc.service.DeleteProductByID(c.Request.Context(), productID)
	if err != nil {
		if errors.Is(err, services.ErrInvalidParams) {
			c.JSON(http.StatusBadRequest,
				apperrors.AppStdErrorHandler(
					services.ErrInvalidParams.Error(),
					"us_0000",
				),
			)
			return
		}
		if errors.Is(err, services.ErrProductDelete) {
			c.JSON(http.StatusUnauthorized,
				apperrors.AppStdErrorHandler(
					services.ErrProductDelete.Error(),
					"us_0001",
				),
			)
			return
		}
		c.JSON(http.StatusInternalServerError,
			apperrors.AppStdErrorHandler(
				"Internal server error",
				"us_0002",
			),
		)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
