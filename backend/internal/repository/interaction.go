package repository

import (
	"context"
	"fmt"

	"jedug-backend/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type InteractionRepo struct {
	pool *pgxpool.Pool
}

func NewInteractionRepo(pool *pgxpool.Pool) *InteractionRepo {
	return &InteractionRepo{pool: pool}
}

// AddReaction adds a reaction (angry, danger, upvote) — 1 per device per report per type
func (r *InteractionRepo) AddReaction(ctx context.Context, input models.CreateInteraction) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO interactions (report_id, fingerprint_hash, type, value)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (report_id, fingerprint_hash, type) DO NOTHING
	`, input.ReportID, input.FingerprintHash, input.Type, input.Value)
	if err != nil {
		return fmt.Errorf("add reaction: %w", err)
	}

	// Update reaction count on report
	_, err = r.pool.Exec(ctx, `
		UPDATE reports SET
			reaction_count = (SELECT COUNT(*) FROM interactions WHERE report_id = $1 AND type != 'comment'),
			updated_at = NOW()
		WHERE id = $1
	`, input.ReportID)
	return err
}

// RemoveReaction removes a specific reaction
func (r *InteractionRepo) RemoveReaction(ctx context.Context, reportID, fingerprint, reactionType string) error {
	_, err := r.pool.Exec(ctx, `
		DELETE FROM interactions
		WHERE report_id = $1 AND fingerprint_hash = $2 AND type = $3
	`, reportID, fingerprint, reactionType)
	if err != nil {
		return fmt.Errorf("remove reaction: %w", err)
	}

	// Update count
	_, err = r.pool.Exec(ctx, `
		UPDATE reports SET
			reaction_count = (SELECT COUNT(*) FROM interactions WHERE report_id = $1 AND type != 'comment'),
			updated_at = NOW()
		WHERE id = $1
	`, reportID)
	return err
}

// AddComment adds a comment to a report
func (r *InteractionRepo) AddComment(ctx context.Context, input models.CreateInteraction) (*models.Interaction, error) {
	var interaction models.Interaction
	err := r.pool.QueryRow(ctx, `
		INSERT INTO interactions (report_id, fingerprint_hash, type, value)
		VALUES ($1, $2, 'comment', $3)
		RETURNING id, report_id, fingerprint_hash, type, value, created_at
	`, input.ReportID, input.FingerprintHash, input.Value).Scan(
		&interaction.ID, &interaction.ReportID, &interaction.FingerprintHash,
		&interaction.Type, &interaction.Value, &interaction.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("add comment: %w", err)
	}
	return &interaction, nil
}

// GetCommentsByReport returns all comments for a report
func (r *InteractionRepo) GetCommentsByReport(ctx context.Context, reportID string) ([]models.Interaction, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, report_id, fingerprint_hash, type, value, created_at
		FROM interactions
		WHERE report_id = $1 AND type = 'comment'
		ORDER BY created_at ASC
	`, reportID)
	if err != nil {
		return nil, fmt.Errorf("get comments: %w", err)
	}
	defer rows.Close()

	var comments []models.Interaction
	for rows.Next() {
		var c models.Interaction
		err := rows.Scan(&c.ID, &c.ReportID, &c.FingerprintHash, &c.Type, &c.Value, &c.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan comment: %w", err)
		}
		comments = append(comments, c)
	}
	return comments, nil
}

// GetReactionsByReport returns all reactions (non-comment) for a report
func (r *InteractionRepo) GetReactionsByReport(ctx context.Context, reportID string) (map[string]int, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT type, COUNT(*) as count
		FROM interactions
		WHERE report_id = $1 AND type != 'comment'
		GROUP BY type
	`, reportID)
	if err != nil {
		return nil, fmt.Errorf("get reactions: %w", err)
	}
	defer rows.Close()

	reactions := make(map[string]int)
	for rows.Next() {
		var t string
		var count int
		if err := rows.Scan(&t, &count); err != nil {
			return nil, fmt.Errorf("scan reaction: %w", err)
		}
		reactions[t] = count
	}
	return reactions, nil
}

// HasReacted checks if a device already reacted with a specific type
func (r *InteractionRepo) HasReacted(ctx context.Context, reportID, fingerprint, reactionType string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM interactions
			WHERE report_id = $1 AND fingerprint_hash = $2 AND type = $3
		)
	`, reportID, fingerprint, reactionType).Scan(&exists)
	return exists, err
}
