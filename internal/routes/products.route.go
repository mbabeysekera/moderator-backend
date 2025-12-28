package routes

import (
	"coolbreez.lk/moderator/internal/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterAdminProductItemsRoutes(routerGroup *gin.RouterGroup,
	authorizationHandler gin.HandlerFunc, rbacHandler gin.HandlerFunc,
	controller *controllers.ProductController) {
	routerGroup.Use(authorizationHandler)
	routerGroup.Use(rbacHandler)
	routerGroup.POST("/products/create", controller.CreateProductWithItems)
	routerGroup.DELETE("/products/delete/:product_id", controller.DeleteProductByID)
}

func RegisterGeneralProductItemsRoutes(routerGroup *gin.RouterGroup,
	controller *controllers.ProductController) {
	routerGroup.GET("/products/all", controller.GetProductsWithItems)
	routerGroup.GET("/products/:product_id", controller.GetProductWithItems)
	routerGroup.GET("/products/sku/:sku", controller.GetProductWithItemsBySku)
	routerGroup.GET("/products/item/:item_code", controller.GetProductWithItemsByItemCode)
}
