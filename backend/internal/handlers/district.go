package handlers

import (
	"github.com/gofiber/fiber/v2"
	"jedug-backend/internal/repository"
)

type DistrictHandler struct {
	repo *repository.DistrictRepo
}

func NewDistrictHandler(repo *repository.DistrictRepo) *DistrictHandler {
	return &DistrictHandler{repo: repo}
}

// GetDistricts returns districts as GeoJSON, optionally filtered by ?city=
func (h *DistrictHandler) GetDistricts(c *fiber.Ctx) error {
	city := c.Query("city", "")

	var data []byte
	var err error

	if city != "" {
		data, err = h.repo.GetDistrictsByCity(c.Context(), city)
	} else {
		data, err = h.repo.GetDistrictsAsGeoJSON(c.Context(), nil)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Set("Content-Type", "application/geo+json")
	return c.Send(data)
}
