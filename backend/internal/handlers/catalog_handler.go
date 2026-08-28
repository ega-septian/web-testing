package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"teststore/backend/internal/models"
)

// CatalogHandler serves read-only, public browse data for the storefront homepage:
// products (New Arrivals / Top Selling) and dress-style categories.
type CatalogHandler struct {
	products    *models.ProductRepo
	dressStyles *models.DressStyleRepo
}

func NewCatalogHandler(products *models.ProductRepo, dressStyles *models.DressStyleRepo) *CatalogHandler {
	return &CatalogHandler{products: products, dressStyles: dressStyles}
}

func (h *CatalogHandler) ListProducts(c *gin.Context) {
	section := c.Query("section")

	products, err := h.products.ListBySection(c.Request.Context(), section)
	if err != nil {
		if errors.Is(err, models.ErrInvalidSection) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parameter 'section' harus 'new_arrivals' atau 'top_selling'"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil daftar produk"})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *CatalogHandler) ListDressStyles(c *gin.Context) {
	styles, err := h.dressStyles.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil daftar dress style"})
		return
	}

	c.JSON(http.StatusOK, styles)
}
