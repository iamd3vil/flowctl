package storage

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// PostgresStorage implements the Storage interface using PostgreSQL
type PostgresStorage struct {
	db *sqlx.DB
}

// NewPostgresStorage creates a new PostgreSQL storage backend
func NewPostgresStorage(db *sqlx.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

// Initialize creates the job queue table
func (p *PostgresStorage) Initialize(ctx context.Context) error {
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS job_queue (
			id SERIAL PRIMARY KEY,
			exec_id TEXT NOT NULL,
			payload JSONB NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);

		CREATE INDEX IF NOT EXISTS idx_job_queue_exec_id ON job_queue(exec_id);
	`

	if _, err := p.db.ExecContext(ctx, createTableQuery); err != nil {
		return err
	}
	return p.migrateAddPayloadType(ctx)
}

// migrateAddPayloadType adds the payload_type column to existing job_queue tables
// This should be called after Initialize
func (p *PostgresStorage) migrateAddPayloadType(ctx context.Context) error {
	// Add payload_type column if it doesn't exist
	addColumnQuery := `
		ALTER TABLE job_queue ADD COLUMN IF NOT EXISTS payload_type TEXT NOT NULL DEFAULT 'flow_execution';
	`
	if _, err := p.db.ExecContext(ctx, addColumnQuery); err != nil {
		return err
	}

	// Create new index for payload_type queries
	createIndexQuery := `
		CREATE INDEX IF NOT EXISTS idx_job_queue_payload_type ON job_queue(payload_type, created_at);
	`
	if _, err := p.db.ExecContext(ctx, createIndexQuery); err != nil {
		return err
	}

	// Drop old index (no longer needed with payload_type index)
	dropOldIndexQuery := `DROP INDEX IF EXISTS idx_job_queue_pending;`
	_, _ = p.db.ExecContext(ctx, dropOldIndexQuery)

	return nil
}

// Put adds a job to the queue
func (p *PostgresStorage) Put(ctx context.Context, job Job) error {
	query := `
		INSERT INTO job_queue (exec_id, payload_type, payload, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := p.db.GetContext(ctx, &job.ID, query, job.ExecID, job.PayloadType, job.Payload, job.CreatedAt)
	return err
}

// GetByPayloadType retrieves and locks a job of specific payload type from the queue
// When the done channel is closed, the job is removed from the queue
func (p *PostgresStorage) GetByPayloadType(ctx context.Context, payloadType string, done chan struct{}) (Job, error) {
	tx, err := p.db.BeginTxx(ctx, nil)
	if err != nil {
		return Job{}, err
	}

	// Select and lock the oldest pending job of this payload type
	selectQuery := `
		SELECT id, exec_id, payload_type, payload, created_at
		FROM job_queue
		WHERE payload_type = $1
		ORDER BY created_at ASC
		LIMIT 1
		FOR UPDATE SKIP LOCKED
	`

	var job Job
	err = tx.GetContext(ctx, &job, selectQuery, payloadType)
	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return Job{}, ErrNoJobs
		}
		return Job{}, err
	}

	// Wait for job completion in background, then delete and commit
	go func() {
		<-done

		deleteQuery := `DELETE FROM job_queue WHERE id = $1`
		_, _ = tx.ExecContext(context.Background(), deleteQuery, job.ID)
		_ = tx.Commit()
	}()

	return job, nil
}

// Delete removes a job from the queue
func (p *PostgresStorage) Delete(ctx context.Context, jobID int64) error {
	query := `DELETE FROM job_queue WHERE id = $1`
	_, err := p.db.ExecContext(ctx, query, jobID)
	return err
}

// CancelByExecID removes all jobs with the given execution ID
func (p *PostgresStorage) CancelByExecID(ctx context.Context, execID string) error {
	query := `DELETE FROM job_queue WHERE exec_id = $1`
	_, err := p.db.ExecContext(ctx, query, execID)
	return err
}

// Close closes the storage backend
func (p *PostgresStorage) Close() error {
	// The database connection is managed externally, so we don't close it here
	return nil
}
