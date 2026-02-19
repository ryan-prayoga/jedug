package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UploadHandler struct {
	uploadDir string
	baseURL   string
}

func NewUploadHandler(uploadDir, baseURL string) *UploadHandler {
	// Create upload directory if not exists
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create upload dir: %v", err))
	}
	return &UploadHandler{
		uploadDir: uploadDir,
		baseURL:   baseURL,
	}
}

// UploadImage handles POST /api/v1/upload
func (h *UploadHandler) UploadImage(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Image file is required",
		})
	}

	// Validate file size (max 5MB — after client compression should be <200KB)
	if file.Size > 5*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File size exceeds 5MB limit",
		})
	}

	// Validate file type
	contentType := file.Header.Get("Content-Type")
	validTypes := map[string]string{
		"image/webp": ".webp",
		"image/jpeg": ".jpg",
		"image/png":  ".png",
	}

	ext, ok := validTypes[contentType]
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Only WebP, JPEG, and PNG images are allowed",
		})
	}

	// Generate unique filename with date-based directory
	now := time.Now()
	subDir := now.Format("2006/01")
	fullDir := filepath.Join(h.uploadDir, subDir)
	if err := os.MkdirAll(fullDir, 0755); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create directory",
		})
	}

	filename := uuid.New().String() + ext
	savePath := filepath.Join(fullDir, filename)

	if err := c.SaveFile(file, savePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save file",
		})
	}

	// Build the URL path
	urlPath := strings.TrimRight(h.baseURL, "/") + "/uploads/" + subDir + "/" + filename

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"url":      urlPath,
		"filename": filename,
		"size":     file.Size,
		"message":  "Upload berhasil",
	})
}
