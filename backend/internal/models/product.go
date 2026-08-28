package models

type Product struct {
	ID          string `json:"id"`
	Brand       string `json:"brand"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Gender      string `json:"gender"`
	Category    string `json:"category"`
	Subcategory string `json:"subcategory"`
	Price       int    `json:"price"`
	// Discount is a percentage (0-100). There's no separate "original price"
	// stored — the UI derives it from Price and Discount when it needs to
	// show a struck-through price.
	Discount  int `json:"discount"`
	TotalSold int `json:"total_sold"`
	// ImageURL is the first uploaded photo (by display_order), or nil if the
	// product has none yet — there's no icon/emoji fallback anymore.
	ImageURL  *string  `json:"image_url"`
	CreatedAt JSONTime `json:"created_at"`
}

// ProductDetail is a Product plus its per-size stock and full gallery
// (ImageURL above is just the first of these), returned by the
// single-product detail endpoint.
type ProductDetail struct {
	Product
	Sizes  []ProductSize  `json:"sizes"`
	Images []ProductImage `json:"images"`
}

// FilterOption is one selectable value in a filter facet, with how many
// products currently match it.
type FilterOption struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

// ProductFilterOptions is the full facet set for the Shop page's sidebar —
// computed live from actual product data, not a hardcoded taxonomy.
type ProductFilterOptions struct {
	Gender      []FilterOption `json:"gender"`
	Category    []FilterOption `json:"category"`
	Subcategory []FilterOption `json:"subcategory"`
	Size        []FilterOption `json:"size"`
}

// ProductFilters narrows ProductRepo.List to products matching all of the
// given (non-empty) fields. Each field is a checklist — multiple values in
// the same dimension are OR'd together (e.g. Gender: ["Pria","Wanita"]
// matches either), while different dimensions are AND'd.
type ProductFilters struct {
	Gender      []string
	Category    []string
	Subcategory []string
	Size        []string
	// Query is free-text search (from the homepage/Shop search box), matched
	// case-insensitively against name, brand, and description. Empty means
	// unfiltered, same as the other fields.
	Query string
}
