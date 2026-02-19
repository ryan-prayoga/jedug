package models

import "time"

type District struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Level string `json:"level"`
}

type Report struct {
	ID              string    `json:"id"`
	FingerprintHash string    `json:"fingerprint_hash,omitempty"`
	UserID          *string   `json:"user_id,omitempty"`
	Latitude        float64   `json:"latitude"`
	Longitude       float64   `json:"longitude"`
	DistrictID      *int      `json:"district_id,omitempty"`
	DistrictName    string    `json:"district_name,omitempty"`
	ImageURL        string    `json:"image_url"`
	Status          string    `json:"status"`
	Severity        int       `json:"severity"`
	RoadType        *string   `json:"road_type,omitempty"`
	ViewCount       int       `json:"view_count"`
	ReactionCount   int       `json:"reaction_count"`
	EstimatedLoss   float64   `json:"estimated_loss"`
	Description     *string   `json:"description,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	// Computed fields for display
	DaysOld int      `json:"days_old,omitempty"`
	Photos  []string `json:"photos,omitempty"`
}

type CreateReportInput struct {
	FingerprintHash string  `json:"fingerprint_hash" validate:"required"`
	UserID          *string `json:"user_id,omitempty"`
	Latitude        float64 `json:"latitude" validate:"required"`
	Longitude       float64 `json:"longitude" validate:"required"`
	ImageURL        string  `json:"image_url" validate:"required"`
	Severity        int     `json:"severity" validate:"required,min=1,max=5"`
	Description     *string `json:"description,omitempty"`
	RoadType        *string `json:"road_type,omitempty"`
}

type ReportFilter struct {
	Severity   string `json:"severity,omitempty"`
	Status     string `json:"status,omitempty"`
	DistrictID int    `json:"district_id,omitempty"`
	Limit      int    `json:"limit,omitempty"`
	Offset     int    `json:"offset,omitempty"`
}

type ReportProof struct {
	ID              string    `json:"id"`
	ReportID        string    `json:"report_id"`
	FingerprintHash string    `json:"fingerprint_hash"`
	ImageURL        string    `json:"image_url"`
	TakenAt         time.Time `json:"taken_at"`
}

type Interaction struct {
	ID              int       `json:"id"`
	ReportID        string    `json:"report_id"`
	FingerprintHash string    `json:"fingerprint_hash"`
	Type            string    `json:"type"`
	Value           *string   `json:"value,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

type CreateInteraction struct {
	ReportID        string  `json:"report_id" validate:"required"`
	FingerprintHash string  `json:"fingerprint_hash" validate:"required"`
	Type            string  `json:"type" validate:"required"`
	Value           *string `json:"value,omitempty"`
}

type User struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	Name       *string   `json:"name,omitempty"`
	AvatarURL  *string   `json:"avatar_url,omitempty"`
	XPPoints   int       `json:"xp_points"`
	RankTitle  string    `json:"rank_title"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
}

type KecamatanRank struct {
	Rank         int     `json:"rank"`
	Name         string  `json:"name"`
	TotalReports int     `json:"total_reports"`
	TotalLoss    float64 `json:"total_loss"`
	PercentFixed int     `json:"percent_fixed"`
	TopSeverity  int     `json:"top_severity"`
}

// API Response types
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int         `json:"total"`
	Limit      int         `json:"limit"`
	Offset     int         `json:"offset"`
	HasMore    bool        `json:"has_more"`
}
