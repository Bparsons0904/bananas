package inserters

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PGXCopyInserter uses PostgreSQL COPY protocol for maximum performance
type PGXCopyInserter struct {
	pool *pgxpool.Pool
}

// NewPGXCopyInserter creates a new PGX COPY inserter
func NewPGXCopyInserter(pool *pgxpool.Pool) *PGXCopyInserter {
	return &PGXCopyInserter{pool: pool}
}

// BulkInsert performs a bulk insert using COPY protocol
func (p *PGXCopyInserter) BulkInsert(ctx context.Context, tableName string, columns []string, rows [][]interface{}) error {
	if len(rows) == 0 {
		return nil
	}

	conn, err := p.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	copyCount, err := conn.Conn().CopyFrom(
		ctx,
		pgx.Identifier{tableName},
		columns,
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return fmt.Errorf("copy failed: %w", err)
	}

	if int(copyCount) != len(rows) {
		return fmt.Errorf("copy count mismatch: expected %d, got %d", len(rows), copyCount)
	}

	return nil
}

// BulkInsertBatched performs a bulk insert in batches
func (p *PGXCopyInserter) BulkInsertBatched(ctx context.Context, tableName string, columns []string, rows [][]interface{}, batchSize int) error {
	if len(rows) == 0 {
		return nil
	}

	for i := 0; i < len(rows); i += batchSize {
		end := i + batchSize
		if end > len(rows) {
			end = len(rows)
		}

		batch := rows[i:end]
		if err := p.BulkInsert(ctx, tableName, columns, batch); err != nil {
			return fmt.Errorf("batch insert failed at offset %d: %w", i, err)
		}
	}

	return nil
}

// InsertReturningIDs inserts rows and returns their generated IDs
// Useful when we need the auto-generated IDs for foreign key relationships
func (p *PGXCopyInserter) InsertReturningIDs(ctx context.Context, tableName string, columns []string, rows [][]interface{}) ([]string, error) {
	if len(rows) == 0 {
		return nil, nil
	}

	conn, err := p.pool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	// Build VALUES clause
	valuePlaceholders := make([]string, len(rows))
	args := make([]interface{}, 0, len(rows)*len(columns))
	placeholderIdx := 1

	for i, row := range rows {
		placeholders := make([]string, len(columns))
		for j := range columns {
			placeholders[j] = fmt.Sprintf("$%d", placeholderIdx)
			placeholderIdx++
			args = append(args, row[j])
		}
		valuePlaceholders[i] = "(" + strings.Join(placeholders, ", ") + ")"
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES %s RETURNING id",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(valuePlaceholders, ", "),
	)

	pgxRows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("insert returning failed: %w", err)
	}
	defer pgxRows.Close()

	ids := make([]string, 0, len(rows))
	for pgxRows.Next() {
		var id string
		if err := pgxRows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan id: %w", err)
		}
		ids = append(ids, id)
	}

	if err := pgxRows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return ids, nil
}
