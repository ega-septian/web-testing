package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"teststore/backend/internal/models"
)

// OrderHandler serves checkout and order history. Every route requires auth
// (see main.go's route wiring) — there's no guest checkout on this site.
type OrderHandler struct {
	orders *models.OrderRepo
}

func NewOrderHandler(orders *models.OrderRepo) *OrderHandler {
	return &OrderHandler{orders: orders}
}

type checkoutItemRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Size      string `json:"size" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

type checkoutRequest struct {
	RecipientName string                `json:"recipient_name" binding:"required"`
	Phone         string                `json:"phone" binding:"required"`
	Address       string                `json:"address" binding:"required"`
	Items         []checkoutItemRequest `json:"items" binding:"required,min=1,dive"`
}

// Checkout places an order from the client's cart. Body:
//
//	{"recipient_name": "...", "phone": "...", "address": "...",
//	 "items": [{"product_id": "...", "size": "M", "quantity": 2}]}
//
// Only product_id/size/quantity are taken from the client — name, brand,
// price, and stock are all resolved and validated server-side (see
// OrderRepo.Place), so a stale or tampered cart can't misprice an order or
// oversell stock.
func (h *OrderHandler) Checkout(c *gin.Context) {
	userID := c.GetString("user_id")

	var req checkoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	items := make([]models.CheckoutItem, 0, len(req.Items))
	for _, it := range req.Items {
		items = append(items, models.CheckoutItem{ProductID: it.ProductID, Size: it.Size, Quantity: it.Quantity})
	}

	order, err := h.orders.Place(c.Request.Context(), userID, req.RecipientName, req.Phone, req.Address, items)
	if err != nil {
		var stockErr *models.InsufficientStockError
		var notFoundErr *models.ProductSizeNotFoundError
		switch {
		case errors.As(err, &stockErr):
			c.JSON(http.StatusBadRequest, gin.H{"error": stockErr.Error()})
		case errors.As(err, &notFoundErr):
			c.JSON(http.StatusBadRequest, gin.H{"error": notFoundErr.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat pesanan"})
		}
		return
	}

	c.JSON(http.StatusCreated, order)
}

// ListOrders returns the logged-in user's order history, newest first.
func (h *OrderHandler) ListOrders(c *gin.Context) {
	userID := c.GetString("user_id")

	orders, err := h.orders.ListByUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil daftar pesanan"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetOrder returns one order belonging to the logged-in user.
func (h *OrderHandler) GetOrder(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")

	order, err := h.orders.GetByID(c.Request.Context(), id, userID)
	if err != nil {
		if errors.Is(err, models.ErrOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "pesanan tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil pesanan"})
		return
	}

	c.JSON(http.StatusOK, order)
}
