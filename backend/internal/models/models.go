package models

import "time"

type District struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Level string `json:"level"`
}

type Report struct {
	ID              string    `json:"id"`
	FingerprintHash string    `json:"fingerprint_hash"`
	UserID          *string   `json:"user_id,omitempty"`
	Latitude        float64   `json:"latitude"`
	Longitude       float64   `json:"longitude"`
	DistrictID      *int      `json:"district_id,omitempty"`
	ImageURL        string    `json:"image_url"`
	Status          string    `json:"status"`
	Severity        int       `json:"severity"`
	RoadType        *string   `json:"road_type,omitempty"`
	ViewCount       int       `json:"view_count"`
	ReactionCount   int       `json:"reaction_count"`
	EstimatedLoss   float64   `json:"estimated_loss"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
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
