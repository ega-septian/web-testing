package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"teststore/backend/internal/models"
)

type AssetHandler struct {
	assets       *models.AssetRepo
	uploadDir    string
	maxSizeBytes int64
}

func NewAssetHandler(assets *models.AssetRepo, uploadDir string, maxSizeBytes int64) *AssetHandler {
	return &AssetHandler{assets: assets, uploadDir: uploadDir, maxSizeBytes: maxSizeBytes}
}

var allowedImageTypes = map[string]string{
	"image/jpeg":    ".jpg",
	"image/png":     ".png",
	"image/webp":    ".webp",
	"image/gif":     ".gif",
	"image/svg+xml": ".svg",
}

// Upload stores or replaces the image identified by "key" (e.g. "hero_banner", "logo").
func (h *AssetHandler) Upload(c *gin.Context) {
	key := strings.TrimSpace(c.PostForm("key"))
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "field 'key' wajib diisi"})
		return
	}

	header, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "field 'image' wajib diisi"})
		return
	}

	if header.Size > h.maxSizeBytes {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": fmt.Sprintf("ukuran file maksimal %d MB", h.maxSizeBytes/1024/1024)})
		return
	}

	file, err := header.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membaca file"})
		return
	}
	defer file.Close()

	// Sniff content type from actual bytes rather than trusting the client-supplied header.
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membaca file"})
		return
	}
	contentType := http.DetectContentType(buf[:n])

	ext, ok := allowedImageTypes[contentType]
	if !ok {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "tipe file harus jpg, png, webp, gif, atau svg"})
		return
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal memproses file"})
		return
	}

	suffix, err := randomHex(8)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat nama file"})
		return
	}
	filename := fmt.Sprintf("%s-%s%s", sanitizeKey(key), suffix, ext)

	if err := os.MkdirAll(h.uploadDir, 0o755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyiapkan folder upload"})
		return
	}

	dstPath := filepath.Join(h.uploadDir, filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyimpan file"})
		return
	}
	defer dst.Close()

	written, err := io.Copy(dst, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyimpan file"})
		return
	}

	prev, prevErr := h.assets.FindByKey(c.Request.Context(), key)

	url := "/uploads/" + filename
	asset, err := h.assets.Upsert(c.Request.Context(), key, filename, url, contentType, written)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyimpan metadata asset"})
		return
	}

	// Clean up the file this key previously pointed to, now that the new one is safely stored.
	if prevErr == nil && prev.Filename != filename {
		_ = os.Remove(filepath.Join(h.uploadDir, prev.Filename))
	}

	c.JSON(http.StatusOK, asset)
}

func (h *AssetHandler) Get(c *gin.Context) {
	key := c.Param("key")

	asset, err := h.assets.FindByKey(c.Request.Context(), key)
	if err != nil {
		if errors.Is(err, models.ErrAssetNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "asset tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil asset"})
		return
	}

	c.JSON(http.StatusOK, asset)
}

func (h *AssetHandler) List(c *gin.Context) {
	assets, err := h.assets.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil daftar asset"})
		return
	}

	c.JSON(http.StatusOK, assets)
}

func (h *AssetHandler) Delete(c *gin.Context) {
	key := c.Param("key")

	filename, err := h.assets.Delete(c.Request.Context(), key)
	if err != nil {
		if errors.Is(err, models.ErrAssetNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "asset tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menghapus asset"})
		return
	}

	_ = os.Remove(filepath.Join(h.uploadDir, filename))

	c.Status(http.StatusNoContent)
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func sanitizeKey(key string) string {
	var b strings.Builder
	for _, r := range key {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '-', r == '_':
			b.WriteRune(r)
		default:
			b.WriteRune('-')
		}
	}
	return b.String()
}
