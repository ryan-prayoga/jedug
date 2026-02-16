package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DistrictRepo struct {
	pool *pgxpool.Pool
}

func NewDistrictRepo(pool *pgxpool.Pool) *DistrictRepo {
	return &DistrictRepo{pool: pool}
}

type DistrictGeoJSON struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Level string `json:"level"`
}

// GetDistrictsAsGeoJSON returns all districts as a GeoJSON FeatureCollection
func (r *DistrictRepo) GetDistrictsAsGeoJSON(ctx context.Context, ids []int) ([]byte, error) {
	query := `
		SELECT json_build_object(
			'type', 'FeatureCollection',
			'features', COALESCE(json_agg(
				json_build_object(
					'type', 'Feature',
					'id', d.id,
					'properties', json_build_object(
						'id', d.id,
						'name', d.name,
						'level', d.level
					),
					'geometry', ST_AsGeoJSON(d.geom)::json
				)
			), '[]'::json)
		)
		FROM districts d
	`

	args := []interface{}{}
	if len(ids) > 0 {
		query += " WHERE d.id = ANY($1)"
		args = append(args, ids)
	}

	var result []byte
	var err error
	if len(args) > 0 {
		err = r.pool.QueryRow(ctx, query, args...).Scan(&result)
	} else {
		err = r.pool.QueryRow(ctx, query).Scan(&result)
	}
	if err != nil {
		return nil, fmt.Errorf("query districts geojson: %w", err)
	}

	return result, nil
}

// GetDistrictsByCity returns districts filtered by known city groupings
func (r *DistrictRepo) GetDistrictsByCity(ctx context.Context, city string) ([]byte, error) {
	var ids []int

	switch city {
	case "kota-tangerang":
		ids = []int{15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27}
	case "kab-tangerang":
		ids = []int{119, 120, 121, 122, 123, 124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143, 144, 145, 146, 147}
	case "tangsel":
		ids = []int{148, 149, 150, 151, 152, 153, 154}
	case "kota-serang":
		ids = []int{9, 10, 11, 12, 13, 14}
	case "kota-cilegon":
		ids = []int{1, 2, 3, 4, 5, 6, 7, 8}
	default:
		return r.GetDistrictsAsGeoJSON(ctx, nil)
	}

	return r.GetDistrictsAsGeoJSON(ctx, ids)
}
