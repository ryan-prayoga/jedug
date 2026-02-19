package repository

import (
	"context"
	"fmt"

	"jedug-backend/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReportRepo struct {
	pool *pgxpool.Pool
}

func NewReportRepo(pool *pgxpool.Pool) *ReportRepo {
	return &ReportRepo{pool: pool}
}

// CreateReport inserts a new report, auto-fills district via PostGIS reverse geocoding,
// and checks for smart grouping (existing report within 10m radius).
func (r *ReportRepo) CreateReport(ctx context.Context, input models.CreateReportInput) (*models.Report, error) {
	// 1. Smart Grouping — check if a report exists within 10 meters
	var existingID *string
	err := r.pool.QueryRow(ctx, `
		SELECT id FROM reports
		WHERE ST_DWithin(
			location,
			ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography,
			10
		) AND status != 'archived'
		LIMIT 1
	`, input.Longitude, input.Latitude).Scan(&existingID)

	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("smart grouping check: %w", err)
	}

	// If there's a nearby report, add as proof to that report instead
	if existingID != nil {
		_, err := r.pool.Exec(ctx, `
			INSERT INTO report_proofs (fingerprint_hash, report_id, image_url)
			VALUES ($1, $2, $3)
		`, input.FingerprintHash, *existingID, input.ImageURL)
		if err != nil {
			return nil, fmt.Errorf("add proof to existing: %w", err)
		}
		return r.GetReportByID(ctx, *existingID)
	}

	// 2. Reverse Geocode — find district via PostGIS ST_Contains
	var districtID *int
	err = r.pool.QueryRow(ctx, `
		SELECT id FROM districts
		WHERE ST_Contains(geom, ST_SetSRID(ST_MakePoint($1, $2), 4326))
		LIMIT 1
	`, input.Longitude, input.Latitude).Scan(&districtID)
	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("reverse geocode: %w", err)
	}

	// 3. Insert the report
	var report models.Report
	err = r.pool.QueryRow(ctx, `
		INSERT INTO reports (fingerprint_hash, user_id, location, district_id, image_url, severity, description, road_type)
		VALUES ($1, $2, ST_SetSRID(ST_MakePoint($3, $4), 4326)::geography, $5, $6, $7, $8, $9)
		RETURNING id, fingerprint_hash, user_id,
			ST_Y(location::geometry) as lat, ST_X(location::geometry) as lng,
			district_id, image_url, status, severity, road_type,
			view_count, reaction_count, estimated_loss, description,
			created_at, updated_at
	`, input.FingerprintHash, input.UserID, input.Longitude, input.Latitude,
		districtID, input.ImageURL, input.Severity, input.Description, input.RoadType,
	).Scan(
		&report.ID, &report.FingerprintHash, &report.UserID,
		&report.Latitude, &report.Longitude,
		&report.DistrictID, &report.ImageURL, &report.Status, &report.Severity, &report.RoadType,
		&report.ViewCount, &report.ReactionCount, &report.EstimatedLoss, &report.Description,
		&report.CreatedAt, &report.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("insert report: %w", err)
	}

	return &report, nil
}

// GetReportByID retrieves a single report by its UUID
func (r *ReportRepo) GetReportByID(ctx context.Context, id string) (*models.Report, error) {
	var report models.Report
	err := r.pool.QueryRow(ctx, `
		SELECT r.id, r.fingerprint_hash, r.user_id,
			ST_Y(r.location::geometry) as lat, ST_X(r.location::geometry) as lng,
			r.district_id, r.image_url, r.status, r.severity, r.road_type,
			r.view_count, r.reaction_count, r.estimated_loss, r.description,
			r.created_at, r.updated_at,
			COALESCE(d.name, '') as district_name
		FROM reports r
		LEFT JOIN districts d ON r.district_id = d.id
		WHERE r.id = $1
	`, id).Scan(
		&report.ID, &report.FingerprintHash, &report.UserID,
		&report.Latitude, &report.Longitude,
		&report.DistrictID, &report.ImageURL, &report.Status, &report.Severity, &report.RoadType,
		&report.ViewCount, &report.ReactionCount, &report.EstimatedLoss, &report.Description,
		&report.CreatedAt, &report.UpdatedAt,
		&report.DistrictName,
	)
	if err != nil {
		return nil, fmt.Errorf("get report: %w", err)
	}
	return &report, nil
}

// ListReports returns paginated reports with optional filters
func (r *ReportRepo) ListReports(ctx context.Context, filter models.ReportFilter) ([]models.Report, int, error) {
	// Count query
	countQuery := `SELECT COUNT(*) FROM reports r WHERE 1=1`
	args := []interface{}{}
	argIdx := 1

	if filter.Severity != "" {
		countQuery += fmt.Sprintf(" AND r.severity = $%d", argIdx)
		args = append(args, filter.Severity)
		argIdx++
	}
	if filter.Status != "" {
		countQuery += fmt.Sprintf(" AND r.status = $%d", argIdx)
		args = append(args, filter.Status)
		argIdx++
	}
	if filter.DistrictID > 0 {
		countQuery += fmt.Sprintf(" AND r.district_id = $%d", argIdx)
		args = append(args, filter.DistrictID)
		argIdx++
	}

	var total int
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count reports: %w", err)
	}

	// Data query
	dataQuery := `
		SELECT r.id, r.fingerprint_hash, r.user_id,
			ST_Y(r.location::geometry) as lat, ST_X(r.location::geometry) as lng,
			r.district_id, r.image_url, r.status, r.severity, r.road_type,
			r.view_count, r.reaction_count, r.estimated_loss, r.description,
			r.created_at, r.updated_at,
			COALESCE(d.name, '') as district_name
		FROM reports r
		LEFT JOIN districts d ON r.district_id = d.id
		WHERE 1=1`

	dataArgs := []interface{}{}
	dataIdx := 1

	if filter.Severity != "" {
		dataQuery += fmt.Sprintf(" AND r.severity = $%d", dataIdx)
		dataArgs = append(dataArgs, filter.Severity)
		dataIdx++
	}
	if filter.Status != "" {
		dataQuery += fmt.Sprintf(" AND r.status = $%d", dataIdx)
		dataArgs = append(dataArgs, filter.Status)
		dataIdx++
	}
	if filter.DistrictID > 0 {
		dataQuery += fmt.Sprintf(" AND r.district_id = $%d", dataIdx)
		dataArgs = append(dataArgs, filter.DistrictID)
		dataIdx++
	}

	dataQuery += " ORDER BY r.created_at DESC"

	if filter.Limit > 0 {
		dataQuery += fmt.Sprintf(" LIMIT $%d", dataIdx)
		dataArgs = append(dataArgs, filter.Limit)
		dataIdx++
	}
	if filter.Offset > 0 {
		dataQuery += fmt.Sprintf(" OFFSET $%d", dataIdx)
		dataArgs = append(dataArgs, filter.Offset)
		dataIdx++
	}

	rows, err := r.pool.Query(ctx, dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("list reports: %w", err)
	}
	defer rows.Close()

	var reports []models.Report
	for rows.Next() {
		var rp models.Report
		err := rows.Scan(
			&rp.ID, &rp.FingerprintHash, &rp.UserID,
			&rp.Latitude, &rp.Longitude,
			&rp.DistrictID, &rp.ImageURL, &rp.Status, &rp.Severity, &rp.RoadType,
			&rp.ViewCount, &rp.ReactionCount, &rp.EstimatedLoss, &rp.Description,
			&rp.CreatedAt, &rp.UpdatedAt,
			&rp.DistrictName,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan report: %w", err)
		}
		reports = append(reports, rp)
	}

	return reports, total, nil
}

// GetReportsNearby returns reports within a given radius (meters) from a point
func (r *ReportRepo) GetReportsNearby(ctx context.Context, lat, lng, radiusMeters float64, limit int) ([]models.Report, error) {
	if limit <= 0 {
		limit = 50
	}

	rows, err := r.pool.Query(ctx, `
		SELECT r.id, r.fingerprint_hash, r.user_id,
			ST_Y(r.location::geometry) as lat, ST_X(r.location::geometry) as lng,
			r.district_id, r.image_url, r.status, r.severity, r.road_type,
			r.view_count, r.reaction_count, r.estimated_loss, r.description,
			r.created_at, r.updated_at,
			COALESCE(d.name, '') as district_name,
			ST_Distance(r.location, ST_SetSRID(ST_MakePoint($2, $1), 4326)::geography) as distance
		FROM reports r
		LEFT JOIN districts d ON r.district_id = d.id
		WHERE ST_DWithin(
			r.location,
			ST_SetSRID(ST_MakePoint($2, $1), 4326)::geography,
			$3
		)
		ORDER BY distance ASC
		LIMIT $4
	`, lat, lng, radiusMeters, limit)
	if err != nil {
		return nil, fmt.Errorf("nearby reports: %w", err)
	}
	defer rows.Close()

	var reports []models.Report
	for rows.Next() {
		var rp models.Report
		var distance float64
		err := rows.Scan(
			&rp.ID, &rp.FingerprintHash, &rp.UserID,
			&rp.Latitude, &rp.Longitude,
			&rp.DistrictID, &rp.ImageURL, &rp.Status, &rp.Severity, &rp.RoadType,
			&rp.ViewCount, &rp.ReactionCount, &rp.EstimatedLoss, &rp.Description,
			&rp.CreatedAt, &rp.UpdatedAt,
			&rp.DistrictName,
			&distance,
		)
		if err != nil {
			return nil, fmt.Errorf("scan nearby: %w", err)
		}
		reports = append(reports, rp)
	}

	return reports, nil
}

// IncrementViewCount adds a view to a report
func (r *ReportRepo) IncrementViewCount(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE reports SET view_count = view_count + 1, updated_at = NOW()
		WHERE id = $1
	`, id)
	return err
}

// UpdateEstimatedLoss recalculates the estimated loss using the viral formula
// Formula: (view_count x 500) + (days_old x 100000)
func (r *ReportRepo) UpdateEstimatedLoss(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE reports SET
			estimated_loss = (view_count * 500) + (EXTRACT(EPOCH FROM (NOW() - created_at)) / 86400 * 100000),
			updated_at = NOW()
		WHERE id = $1
	`, id)
	return err
}

// GetProofsByReportID returns all additional proofs for a report
func (r *ReportRepo) GetProofsByReportID(ctx context.Context, reportID string) ([]models.ReportProof, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, report_id, fingerprint_hash, image_url, taken_at
		FROM report_proofs
		WHERE report_id = $1
		ORDER BY taken_at DESC
	`, reportID)
	if err != nil {
		return nil, fmt.Errorf("get proofs: %w", err)
	}
	defer rows.Close()

	var proofs []models.ReportProof
	for rows.Next() {
		var p models.ReportProof
		err := rows.Scan(&p.ID, &p.ReportID, &p.FingerprintHash, &p.ImageURL, &p.TakenAt)
		if err != nil {
			return nil, fmt.Errorf("scan proof: %w", err)
		}
		proofs = append(proofs, p)
	}
	return proofs, nil
}

// GetKecamatanRanking returns the top kecamatan ranked by total reports + loss
func (r *ReportRepo) GetKecamatanRanking(ctx context.Context, limit int) ([]models.KecamatanRank, error) {
	if limit <= 0 {
		limit = 10
	}

	rows, err := r.pool.Query(ctx, `
		SELECT
			d.name,
			COUNT(r.id) as total_reports,
			COALESCE(SUM(r.estimated_loss), 0) as total_loss,
			COALESCE(
				ROUND(COUNT(CASE WHEN r.status = 'fixed' THEN 1 END)::numeric / NULLIF(COUNT(r.id), 0) * 100),
				0
			) as percent_fixed,
			COALESCE(
				(SELECT severity FROM reports r2
				 WHERE r2.district_id = d.id AND r2.status != 'archived'
				 ORDER BY r2.severity DESC LIMIT 1),
				1
			) as top_severity
		FROM districts d
		LEFT JOIN reports r ON r.district_id = d.id AND r.status != 'archived'
		WHERE d.level = 'kecamatan'
		GROUP BY d.id, d.name
		HAVING COUNT(r.id) > 0
		ORDER BY total_reports DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("kecamatan ranking: %w", err)
	}
	defer rows.Close()

	var ranks []models.KecamatanRank
	for rows.Next() {
		var rank models.KecamatanRank
		err := rows.Scan(
			&rank.Name, &rank.TotalReports, &rank.TotalLoss,
			&rank.PercentFixed, &rank.TopSeverity,
		)
		if err != nil {
			return nil, fmt.Errorf("scan rank: %w", err)
		}
		ranks = append(ranks, rank)
	}

	// Assign ranks
	for i := range ranks {
		ranks[i].Rank = i + 1
	}

	return ranks, nil
}
