package models

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrOrderNotFound = errors.New("order not found")

// InsufficientStockError reports that a specific product+size in the
// checkout request doesn't have enough stock left. Surfaced as a 400, not a
// 500 — running out of stock mid-checkout (e.g. two people buying the last
// unit at once) is a normal, expected outcome, not a server bug.
type InsufficientStockError struct {
	ProductName string
	Size        string
	Available   int
}

func (e *InsufficientStockError) Error() string {
	return fmt.Sprintf("stok tidak cukup untuk %s ukuran %s (tersisa %d)", e.ProductName, e.Size, e.Available)
}

// ProductSizeNotFoundError reports a checkout line referencing a product+size
// combination that doesn't exist (bad/stale product ID, or a size the
// product doesn't offer) — a client-data problem, not a server bug.
type ProductSizeNotFoundError struct {
	ProductID string
	Size      string
}

func (e *ProductSizeNotFoundError) Error() string {
	return fmt.Sprintf("produk %s ukuran %s tidak ditemukan", e.ProductID, e.Size)
}

type OrderRepo struct {
	pool *pgxpool.Pool
}

func NewOrderRepo(pool *pgxpool.Pool) *OrderRepo {
	return &OrderRepo{pool: pool}
}

// Place validates + locks stock, decrements it, and creates the order plus
// its items and a matching sales row per item (so best_selling stays
// accurate now that real orders exist), all in one transaction — either the
// whole checkout succeeds, or none of it does. Price is always the product's
// current price, looked up here — never trusted from the client, so a
// stale/tampered cart can't under- or over-charge.
func (r *OrderRepo) Place(
	ctx context.Context, userID, recipientName, phone, address string, requested []CheckoutItem,
) (*Order, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck // no-op if Commit already succeeded

	resolved := make([]OrderItem, 0, len(requested))
	totalAmount := 0

	for _, req := range requested {
		var productName, brand string
		var unitPrice, stock int
		// FOR UPDATE locks the row so two concurrent checkouts on the same
		// product+size can't both read "stock available" before either one
		// decrements it.
		err := tx.QueryRow(ctx,
			`SELECT p.name, p.brand, p.price, ps.stock
			 FROM product_sizes ps
			 JOIN products p ON p.id = ps.product_id
			 WHERE ps.product_id = $1 AND ps.size = $2
			 FOR UPDATE OF ps`,
			req.ProductID, req.Size,
		).Scan(&productName, &brand, &unitPrice, &stock)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &ProductSizeNotFoundError{ProductID: req.ProductID, Size: req.Size}
		}
		if err != nil {
			return nil, err
		}
		if stock < req.Quantity {
			return nil, &InsufficientStockError{ProductName: productName, Size: req.Size, Available: stock}
		}

		productID := req.ProductID
		resolved = append(resolved, OrderItem{
			ProductID: &productID, ProductName: productName, Brand: brand,
			Size: req.Size, UnitPrice: unitPrice, Quantity: req.Quantity,
		})
		totalAmount += unitPrice * req.Quantity
	}

	var order Order
	err = tx.QueryRow(ctx,
		`INSERT INTO orders (user_id, recipient_name, phone, address, total_amount, status)
		 VALUES ($1, $2, $3, $4, $5, 'placed')
		 RETURNING id, user_id, recipient_name, phone, address, total_amount, status, created_at`,
		userID, recipientName, phone, address, totalAmount,
	).Scan(
		&order.ID, &order.UserID, &order.RecipientName, &order.Phone, &order.Address,
		&order.TotalAmount, &order.Status, &order.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	for i := range resolved {
		item := &resolved[i]

		if _, err := tx.Exec(ctx,
			`UPDATE product_sizes SET stock = stock - $1 WHERE product_id = $2 AND size = $3`,
			item.Quantity, *item.ProductID, item.Size,
		); err != nil {
			return nil, err
		}

		if err := tx.QueryRow(ctx,
			`INSERT INTO order_items (order_id, product_id, product_name, brand, size, unit_price, quantity)
			 VALUES ($1, $2, $3, $4, $5, $6, $7)
			 RETURNING id, created_at`,
			order.ID, item.ProductID, item.ProductName, item.Brand, item.Size, item.UnitPrice, item.Quantity,
		).Scan(&item.ID, &item.CreatedAt); err != nil {
			return nil, err
		}

		if _, err := tx.Exec(ctx,
			`INSERT INTO sales (product_id, quantity) VALUES ($1, $2)`,
			*item.ProductID, item.Quantity,
		); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	order.Items = resolved
	return &order, nil
}

// ListByUser returns the given user's orders, newest first, each with its items.
func (r *OrderRepo) ListByUser(ctx context.Context, userID string) ([]Order, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, recipient_name, phone, address, total_amount, status, created_at
		 FROM orders WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []Order{}
	for rows.Next() {
		var o Order
		if err := rows.Scan(
			&o.ID, &o.UserID, &o.RecipientName, &o.Phone, &o.Address, &o.TotalAmount, &o.Status, &o.CreatedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// N+1 (one items query per order) — fine at this catalog's scale, same
	// trade-off already made elsewhere (e.g. GetProduct's sizes+images).
	for i := range orders {
		items, err := r.itemsByOrder(ctx, orders[i].ID)
		if err != nil {
			return nil, err
		}
		orders[i].Items = items
	}

	return orders, nil
}

// GetByID returns one order, scoped to userID — an order that exists but
// belongs to someone else is indistinguishable from a nonexistent one.
func (r *OrderRepo) GetByID(ctx context.Context, id, userID string) (*Order, error) {
	var o Order
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, recipient_name, phone, address, total_amount, status, created_at
		 FROM orders WHERE id = $1 AND user_id = $2`,
		id, userID,
	).Scan(&o.ID, &o.UserID, &o.RecipientName, &o.Phone, &o.Address, &o.TotalAmount, &o.Status, &o.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrOrderNotFound
	}
	if err != nil {
		return nil, err
	}

	items, err := r.itemsByOrder(ctx, o.ID)
	if err != nil {
		return nil, err
	}
	o.Items = items

	return &o, nil
}

func (r *OrderRepo) itemsByOrder(ctx context.Context, orderID string) ([]OrderItem, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, product_id, product_name, brand, size, unit_price, quantity, created_at
		 FROM order_items WHERE order_id = $1 ORDER BY created_at ASC`,
		orderID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []OrderItem{}
	for rows.Next() {
		var it OrderItem
		if err := rows.Scan(
			&it.ID, &it.ProductID, &it.ProductName, &it.Brand, &it.Size, &it.UnitPrice, &it.Quantity, &it.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, it)
	}

	return items, rows.Err()
}
