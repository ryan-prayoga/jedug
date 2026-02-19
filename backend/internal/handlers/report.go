package handlers

import (
	"math"
	"strconv"
	"time"

	"jedug-backend/internal/models"
	"jedug-backend/internal/repository"

	"github.com/gofiber/fiber/v2"
)

type ReportHandler struct {
	reportRepo      *repository.ReportRepo
	interactionRepo *repository.InteractionRepo
}

func NewReportHandler(rr *repository.ReportRepo, ir *repository.InteractionRepo) *ReportHandler {
	return &ReportHandler{reportRepo: rr, interactionRepo: ir}
}

// CreateReport handles POST /api/v1/reports
func (h *ReportHandler) CreateReport(c *fiber.Ctx) error {
	var input models.CreateReportInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validation
	if input.FingerprintHash == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "fingerprint_hash is required",
		})
	}
	if input.Latitude == 0 || input.Longitude == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "latitude and longitude are required",
		})
	}
	if input.ImageURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "image_url is required",
		})
	}
	if input.Severity < 1 || input.Severity > 5 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "severity must be between 1 and 5",
		})
	}

	report, err := h.reportRepo.CreateReport(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    report,
		"message": "Laporan berhasil dikirim!",
	})
}

// GetReport handles GET /api/v1/reports/:id
func (h *ReportHandler) GetReport(c *fiber.Ctx) error {
	id := c.Params("id")

	report, err := h.reportRepo.GetReportByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Laporan tidak ditemukan",
		})
	}

	// Compute dynamic fields
	report.DaysOld = int(math.Floor(time.Since(report.CreatedAt).Hours() / 24))

	// Get proofs (additional photos)
	proofs, _ := h.reportRepo.GetProofsByReportID(c.Context(), id)
	photos := []string{report.ImageURL}
	for _, p := range proofs {
		photos = append(photos, p.ImageURL)
	}
	report.Photos = photos

	// Increment view count (fire & forget)
	go func() {
		_ = h.reportRepo.IncrementViewCount(c.Context(), id)
		_ = h.reportRepo.UpdateEstimatedLoss(c.Context(), id)
	}()

	return c.JSON(fiber.Map{"data": report})
}

// ListReports handles GET /api/v1/reports
func (h *ReportHandler) ListReports(c *fiber.Ctx) error {
	filter := models.ReportFilter{
		Severity: c.Query("severity", ""),
		Status:   c.Query("status", ""),
	}

	if did := c.Query("district_id", ""); did != "" {
		if v, err := strconv.Atoi(did); err == nil {
			filter.DistrictID = v
		}
	}

	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > 100 {
		limit = 100
	}
	filter.Limit = limit
	filter.Offset = offset

	reports, total, err := h.reportRepo.ListReports(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Compute dynamic fields
	now := time.Now()
	for i := range reports {
		reports[i].DaysOld = int(math.Floor(now.Sub(reports[i].CreatedAt).Hours() / 24))
		reports[i].Photos = []string{reports[i].ImageURL}
	}

	return c.JSON(models.PaginatedResponse{
		Data:    reports,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: offset+limit < total,
	})
}

// GetNearbyReports handles GET /api/v1/reports/nearby
func (h *ReportHandler) GetNearbyReports(c *fiber.Ctx) error {
	lat, err := strconv.ParseFloat(c.Query("lat", "0"), 64)
	if err != nil || lat == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "lat is required",
		})
	}
	lng, err := strconv.ParseFloat(c.Query("lng", "0"), 64)
	if err != nil || lng == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "lng is required",
		})
	}

	radius, _ := strconv.ParseFloat(c.Query("radius", "1000"), 64)
	limit, _ := strconv.Atoi(c.Query("limit", "50"))

	reports, err := h.reportRepo.GetReportsNearby(c.Context(), lat, lng, radius, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{"data": reports})
}

// AddReaction handles POST /api/v1/reports/:id/reactions
func (h *ReportHandler) AddReaction(c *fiber.Ctx) error {
	reportID := c.Params("id")

	var input struct {
		FingerprintHash string `json:"fingerprint_hash"`
		Type            string `json:"type"` // angry, danger, upvote
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if input.FingerprintHash == "" || input.Type == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "fingerprint_hash and type are required",
		})
	}

	// Valid reaction types
	validTypes := map[string]bool{"angry": true, "danger": true, "upvote": true}
	if !validTypes[input.Type] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid reaction type. Use: angry, danger, upvote",
		})
	}

	err := h.interactionRepo.AddReaction(c.Context(), models.CreateInteraction{
		ReportID:        reportID,
		FingerprintHash: input.FingerprintHash,
		Type:            input.Type,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{"message": "Reaksi ditambahkan"})
}

// AddComment handles POST /api/v1/reports/:id/comments
func (h *ReportHandler) AddComment(c *fiber.Ctx) error {
	reportID := c.Params("id")

	var input struct {
		FingerprintHash string `json:"fingerprint_hash"`
		Text            string `json:"text"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if input.FingerprintHash == "" || input.Text == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "fingerprint_hash and text are required",
		})
	}

	comment, err := h.interactionRepo.AddComment(c.Context(), models.CreateInteraction{
		ReportID:        reportID,
		FingerprintHash: input.FingerprintHash,
		Type:            "comment",
		Value:           &input.Text,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": comment})
}

// GetComments handles GET /api/v1/reports/:id/comments
func (h *ReportHandler) GetComments(c *fiber.Ctx) error {
	reportID := c.Params("id")

	comments, err := h.interactionRepo.GetCommentsByReport(c.Context(), reportID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{"data": comments})
}

// GetRanking handles GET /api/v1/ranking
func (h *ReportHandler) GetRanking(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	ranks, err := h.reportRepo.GetKecamatanRanking(c.Context(), limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{"data": ranks})
}
