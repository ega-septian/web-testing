package models

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrInvalidSort = errors.New("invalid product sort")
var ErrProductNotFound = errors.New("product not found")

// validSorts maps a public ?sort= value to a safe, hardcoded ORDER BY
// fragment — never built from raw user input, so this stays injection-safe
// despite being interpolated into the query string.
var validSorts = map[string]string{
	"newest":       "p.created_at DESC",
	"best_selling": "total_sold DESC, p.created_at DESC",
}

const productSelectColumns = `
	p.id, p.brand, p.name, p.description, p.gender, p.category, p.subcategory, p.price, p.discount, p.created_at,
	COALESCE(SUM(s.quantity), 0) AS total_sold,
	(SELECT pi.url FROM product_images pi WHERE pi.product_id = p.id ORDER BY pi.display_order ASC LIMIT 1) AS image_url`

type ProductRepo struct {
	pool *pgxpool.Pool
}

func NewProductRepo(pool *pgxpool.Pool) *ProductRepo {
	return &ProductRepo{pool: pool}
}

func scanProduct(row pgx.Row) (*Product, error) {
	var p Product
	err := row.Scan(
		&p.ID, &p.Brand, &p.Name, &p.Description, &p.Gender, &p.Category, &p.Subcategory,
		&p.Price, &p.Discount, &p.CreatedAt, &p.TotalSold, &p.ImageURL,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// List returns products across the whole catalog, ranked by "newest"
// (created_at) or "best_selling" (total units sold, from the sales table),
// optionally narrowed by filters (all fields optional — empty means
// unfiltered on that dimension).
func (r *ProductRepo) List(ctx context.Context, sort string, limit int, filters ProductFilters) ([]Product, error) {
	orderBy, ok := validSorts[sort]
	if !ok {
		return nil, ErrInvalidSort
	}

	var conditions []string
	var args []any

	// Multiple values within one dimension are OR'd (checklist semantics):
	// gender IN ('Pria','Wanita') matches either. Different dimensions AND.
	addAny := func(column string, values []string) {
		if len(values) == 0 {
			return
		}
		args = append(args, values)
		conditions = append(conditions, fmt.Sprintf("%s = ANY($%d)", column, len(args)))
	}
	addAny("p.brand", filters.Brand)
	addAny("p.gender", filters.Gender)
	addAny("p.category", filters.Category)
	addAny("p.subcategory", filters.Subcategory)

	if len(filters.Size) > 0 {
		args = append(args, filters.Size)
		conditions = append(conditions, fmt.Sprintf(
			"EXISTS (SELECT 1 FROM product_sizes ps WHERE ps.product_id = p.id AND ps.size = ANY($%d))", len(args),
		))
	}

	// Free-text search matches name, brand, or description — same $N
	// placeholder reused three times, so it only costs one bound arg.
	if q := strings.TrimSpace(filters.Query); q != "" {
		args = append(args, "%"+q+"%")
		idx := len(args)
		conditions = append(conditions, fmt.Sprintf(
			"(p.name ILIKE $%d OR p.brand ILIKE $%d OR p.description ILIKE $%d)", idx, idx, idx,
		))
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	args = append(args, limit)
	query := fmt.Sprintf(`
		SELECT %s
		FROM products p
		LEFT JOIN sales s ON s.product_id = p.id
		%s
		GROUP BY p.id
		ORDER BY %s
		LIMIT $%d`, productSelectColumns, where, orderBy, len(args))

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		p, err := scanProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, *p)
	}

	return products, rows.Err()
}

// GetByID returns a single product (with total_sold and its primary image),
// or ErrProductNotFound.
func (r *ProductRepo) GetByID(ctx context.Context, id string) (*Product, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM products p
		LEFT JOIN sales s ON s.product_id = p.id
		WHERE p.id = $1
		GROUP BY p.id`, productSelectColumns)

	p, err := scanProduct(r.pool.QueryRow(ctx, query, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrProductNotFound
	}
	if err != nil {
		return nil, err
	}

	return p, nil
}

// filterFacet counts distinct values for one text column on products,
// skipping the empty string (products that were never given a value for it).
func (r *ProductRepo) filterFacet(ctx context.Context, column string) ([]FilterOption, error) {
	rows, err := r.pool.Query(ctx, fmt.Sprintf(
		`SELECT %s, COUNT(*) FROM products WHERE %s != '' GROUP BY %s ORDER BY %s ASC`, column, column, column, column,
	))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	options := []FilterOption{}
	for rows.Next() {
		var o FilterOption
		if err := rows.Scan(&o.Value, &o.Count); err != nil {
			return nil, err
		}
		options = append(options, o)
	}

	return options, rows.Err()
}

// FilterOptions computes the Shop page's sidebar facets from live product
// data — counts are global (not narrowed by other currently-applied
// filters); see the handler doc comment for that trade-off.
func (r *ProductRepo) FilterOptions(ctx context.Context) (*ProductFilterOptions, error) {
	brand, err := r.filterFacet(ctx, "brand")
	if err != nil {
		return nil, err
	}
	gender, err := r.filterFacet(ctx, "gender")
	if err != nil {
		return nil, err
	}
	category, err := r.filterFacet(ctx, "category")
	if err != nil {
		return nil, err
	}
	subcategory, err := r.filterFacet(ctx, "subcategory")
	if err != nil {
		return nil, err
	}

	sizeRows, err := r.pool.Query(ctx,
		`SELECT size, COUNT(DISTINCT product_id) FROM product_sizes GROUP BY size ORDER BY size ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer sizeRows.Close()

	size := []FilterOption{}
	for sizeRows.Next() {
		var o FilterOption
		if err := sizeRows.Scan(&o.Value, &o.Count); err != nil {
			return nil, err
		}
		size = append(size, o)
	}
	if err := sizeRows.Err(); err != nil {
		return nil, err
	}

	return &ProductFilterOptions{Brand: brand, Gender: gender, Category: category, Subcategory: subcategory, Size: size}, nil
}

// Create adds a new product plus its per-size stock, in one transaction —
// either both are written, or neither is. Photos are added afterward in the
// same request by the handler (see CatalogHandler.CreateProduct), so
// ImageURL is always nil on the object this returns.
func (r *ProductRepo) Create(
	ctx context.Context, brand, name, description, gender, category, subcategory string,
	price, discount int, sizes []ProductSize,
) (*Product, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck // no-op if Commit already succeeded

	var p Product
	err = tx.QueryRow(ctx,
		`INSERT INTO products (brand, name, description, gender, category, subcategory, price, discount)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id, brand, name, description, gender, category, subcategory, price, discount, created_at`,
		brand, name, description, gender, category, subcategory, price, discount,
	).Scan(
		&p.ID, &p.Brand, &p.Name, &p.Description, &p.Gender, &p.Category, &p.Subcategory,
		&p.Price, &p.Discount, &p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	for _, s := range sizes {
		if _, err := tx.Exec(ctx,
			`INSERT INTO product_sizes (product_id, size, stock) VALUES ($1, $2, $3)`,
			p.ID, s.Size, s.Stock,
		); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &p, nil
}
