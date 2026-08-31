package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"teststore/backend/internal/models"
)

// DefaultProductsLimit caps a product listing when the caller doesn't
// override it via ?limit=.
const DefaultProductsLimit = 4

// validBrands mirrors the brand strip on the homepage (HomeView.vue) — a
// product's brand must be one of these, not arbitrary free text.
var validBrands = map[string]bool{
	"NEVADA": true,
	"DISNEY": true,
	"MARVEL": true,
	"COLE":   true,
	"SUKO":   true,
}

// CatalogHandler serves browse data for the storefront homepage: products
// (ranked by recency or sales) and their detail pages. Listing/detail is
// public; creating a product or recording a sale requires auth (see
// main.go's route wiring).
type CatalogHandler struct {
	products     *models.ProductRepo
	productSizes *models.ProductSizeRepo
	productImgs  *models.ProductImageRepo
	sales        *models.SaleRepo
	uploadDir    string
	maxSizeBytes int64
}

func NewCatalogHandler(
	products *models.ProductRepo, productSizes *models.ProductSizeRepo, productImgs *models.ProductImageRepo,
	sales *models.SaleRepo, uploadDir string, maxSizeBytes int64,
) *CatalogHandler {
	return &CatalogHandler{
		products: products, productSizes: productSizes, productImgs: productImgs, sales: sales,
		uploadDir: uploadDir, maxSizeBytes: maxSizeBytes,
	}
}

type recordSaleRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

// ListProducts returns products ranked by ?sort= ("newest" or
// "best_selling"), defaulting to "newest", capped at ?limit= (default
// DefaultProductsLimit). Optionally narrowed by ?brand=, ?gender=,
// ?category=, ?subcategory=, ?size= (all combinable, all optional), and by
// ?q= — a free-text search against name/brand/description, used by the
// homepage search box (which sends the user here, to the Shop page).
func (h *CatalogHandler) ListProducts(c *gin.Context) {
	sort := c.DefaultQuery("sort", "newest")

	limit := DefaultProductsLimit
	if raw := c.Query("limit"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parameter 'limit' harus bilangan bulat positif", "param": "limit"})
			return
		}
		limit = parsed
	}

	// Repeatable query params for checklist-style multi-select, e.g.
	// ?gender=Pria&gender=Wanita.
	filters := models.ProductFilters{
		Brand:       c.QueryArray("brand"),
		Gender:      c.QueryArray("gender"),
		Category:    c.QueryArray("category"),
		Subcategory: c.QueryArray("subcategory"),
		Size:        c.QueryArray("size"),
		Query:       strings.TrimSpace(c.Query("q")),
	}

	products, err := h.products.List(c.Request.Context(), sort, limit, filters)
	if err != nil {
		if errors.Is(err, models.ErrInvalidSort) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parameter 'sort' harus 'newest' atau 'best_selling'", "param": "sort"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil daftar produk"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// ListProductFilters returns the Shop page's sidebar facets (Brand, Gender,
// Kategori, Sub Kategori, Ukuran) with counts, computed from live product
// data. Simplification worth knowing: each facet's count is global — it
// isn't narrowed by whatever other filters are currently applied, unlike
// e.g. Matahari's live-updating counts. Good enough for this catalog's size;
// revisit if that mismatch becomes noticeable.
func (h *CatalogHandler) ListProductFilters(c *gin.Context) {
	options, err := h.products.FilterOptions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil opsi filter"})
		return
	}

	c.JSON(http.StatusOK, options)
}

// GetProduct returns a single product's detail: per-size stock and gallery images.
func (h *CatalogHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")

	product, err := h.products.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "produk tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil produk"})
		return
	}

	sizes, err := h.productSizes.ListByProduct(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil ukuran produk"})
		return
	}

	images, err := h.productImgs.ListByProduct(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil gambar produk"})
		return
	}

	c.JSON(http.StatusOK, models.ProductDetail{Product: *product, Sizes: sizes, Images: images})
}

// validatedImage is a photo that's passed content-type/size checks, staged to
// be written to disk only after the product itself is created successfully.
type validatedImage struct {
	header *multipart.FileHeader
	ext    string
}

// validateProductImage sniffs the real content type (never trusting the
// client-supplied header) and checks size, without writing anything yet.
func (h *CatalogHandler) validateProductImage(fh *multipart.FileHeader) (validatedImage, error) {
	if fh.Size > h.maxSizeBytes {
		return validatedImage{}, fmt.Errorf("ukuran file '%s' melebihi %d MB", fh.Filename, h.maxSizeBytes/1024/1024)
	}

	f, err := fh.Open()
	if err != nil {
		return validatedImage{}, fmt.Errorf("gagal membaca file '%s'", fh.Filename)
	}
	defer f.Close()

	buf := make([]byte, 512)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return validatedImage{}, fmt.Errorf("gagal membaca file '%s'", fh.Filename)
	}

	ext, ok := allowedImageTypes[http.DetectContentType(buf[:n])]
	if !ok {
		return validatedImage{}, fmt.Errorf("tipe file '%s' harus jpg, png, webp, gif, atau svg", fh.Filename)
	}

	return validatedImage{header: fh, ext: ext}, nil
}

// CreateProduct creates a product plus exactly models.MaxProductImages
// photos in one request. multipart/form-data, not JSON, since it carries files:
//   - brand (must be one of validBrands), name, description, gender, category, subcategory, price, discount: plain form fields
//   - sizes: form field containing a JSON array, e.g. `[{"size":"S","stock":10}]`
//   - images: exactly MaxProductImages file parts (repeat the "images" field per file) — required, not optional
func (h *CatalogHandler) CreateProduct(c *gin.Context) {
	name := strings.TrimSpace(c.PostForm("name"))
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "field 'name' wajib diisi"})
		return
	}

	brand := strings.ToUpper(strings.TrimSpace(c.PostForm("brand")))
	if !validBrands[brand] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "field 'brand' harus salah satu dari: NEVADA, DISNEY, MARVEL, COLE, SUKO"})
		return
	}

	gender := strings.TrimSpace(c.PostForm("gender"))
	if gender == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "field 'gender' wajib diisi"})
		return
	}

	category := strings.TrimSpace(c.PostForm("category"))
	if category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "field 'category' wajib diisi"})
		return
	}

	// Not every category needs a subcategory (matches the reference site,
	// where several categories have none) — optional.
	subcategory := strings.TrimSpace(c.PostForm("subcategory"))

	price, err := strconv.Atoi(c.PostForm("price"))
	if err != nil || price < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "field 'price' wajib diisi, harus bilangan bulat >= 0"})
		return
	}

	discount, err := strconv.Atoi(c.DefaultPostForm("discount", "0"))
	if err != nil || discount < 0 || discount > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "field 'discount' harus bilangan bulat 0-100"})
		return
	}

	var sizes []models.ProductSize
	if raw := c.PostForm("sizes"); raw != "" {
		if err := json.Unmarshal([]byte(raw), &sizes); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": `field 'sizes' harus JSON array, contoh: [{"size":"S","stock":10}]`})
			return
		}
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "gagal membaca form"})
		return
	}
	files := form.File["images"]
	if len(files) != models.MaxProductImages {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("produk wajib punya tepat %d gambar (dapat %d)", models.MaxProductImages, len(files))})
		return
	}

	// Validate every image up front — a bad file should fail before the
	// product is ever created, not leave one half-created behind.
	validated := make([]validatedImage, 0, len(files))
	for _, fh := range files {
		v, err := h.validateProductImage(fh)
		if err != nil {
			c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": err.Error()})
			return
		}
		validated = append(validated, v)
	}

	product, err := h.products.Create(
		c.Request.Context(), brand, name, c.PostForm("description"), gender, category, subcategory, price, discount, sizes,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat produk"})
		return
	}

	// Images are saved after the product exists (they reference its ID) —
	// the product itself is already committed at this point either way, so a
	// failure here is reported but doesn't roll the product back.
	if len(validated) > 0 {
		if err := os.MkdirAll(h.uploadDir, 0o755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "produk berhasil dibuat, tapi gagal menyiapkan folder upload gambar",
				"product": product,
			})
			return
		}

		for i, v := range validated {
			suffix, err := randomHex(8)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "produk berhasil dibuat, tapi gagal menyimpan salah satu gambar",
					"product": product,
				})
				return
			}
			filename := fmt.Sprintf("product-%s-%s%s", sanitizeKey(product.ID), suffix, v.ext)

			src, err := v.header.Open()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "produk berhasil dibuat, tapi gagal menyimpan salah satu gambar",
					"product": product,
				})
				return
			}
			dst, err := os.Create(filepath.Join(h.uploadDir, filename))
			if err != nil {
				src.Close()
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "produk berhasil dibuat, tapi gagal menyimpan salah satu gambar",
					"product": product,
				})
				return
			}
			_, copyErr := io.Copy(dst, src)
			src.Close()
			dst.Close()
			if copyErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "produk berhasil dibuat, tapi gagal menyimpan salah satu gambar",
					"product": product,
				})
				return
			}

			savedImage, err := h.productImgs.Add(c.Request.Context(), product.ID, "/uploads/"+filename)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "produk berhasil dibuat, tapi gagal menyimpan metadata salah satu gambar",
					"product": product,
				})
				return
			}
			// Reflect the upload in this same response — no need for a
			// follow-up GET just to see the photo that was just attached.
			if i == 0 {
				product.ImageURL = &savedImage.URL
			}
		}
	}

	c.JSON(http.StatusCreated, product)
}

// DeleteProduct soft-deletes a product (sets deleted_at instead of removing
// the row — see ProductRepo.SoftDelete) so it disappears from every public
// listing/detail/filter immediately, without touching its size/image/sale
// rows or any order_items that reference it. Meant for cleaning up
// test-seeded data — the storefront itself never calls this.
func (h *CatalogHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	if err := h.products.SoftDelete(c.Request.Context(), id); err != nil {
		if errors.Is(err, models.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "produk tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menghapus produk"})
		return
	}

	c.Status(http.StatusNoContent)
}

// RecordSale logs a sale event for a product, feeding the "best_selling"
// sort. There's no real checkout flow to call this automatically yet — it's
// meant to be called directly (manually, or by whatever replaces this later).
func (h *CatalogHandler) RecordSale(c *gin.Context) {
	productID := c.Param("id")

	if _, err := h.products.GetByID(c.Request.Context(), productID); err != nil {
		if errors.Is(err, models.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "produk tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil produk"})
		return
	}

	var req recordSaleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sale, err := h.sales.Record(c.Request.Context(), productID, req.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mencatat penjualan"})
		return
	}

	c.JSON(http.StatusCreated, sale)
}
