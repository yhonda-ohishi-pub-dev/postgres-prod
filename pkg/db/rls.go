package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RLSContextKey is the context key for RLS organization ID
type rlsContextKey struct{}

// WithOrganizationID returns a new context with organization ID for RLS
func WithOrganizationID(ctx context.Context, orgID string) context.Context {
	return context.WithValue(ctx, rlsContextKey{}, orgID)
}

// GetOrganizationID retrieves the organization ID from context
func GetOrganizationID(ctx context.Context) (string, bool) {
	orgID, ok := ctx.Value(rlsContextKey{}).(string)
	return orgID, ok
}

// SetRLSContext sets the app.current_organization_id session variable for RLS
func SetRLSContext(ctx context.Context, conn *pgxpool.Conn, orgID string) error {
	// SET doesn't support parameterized queries, so we use set_config() function instead
	_, err := conn.Exec(ctx, "SELECT set_config('app.current_organization_id', $1, false)", orgID)
	if err != nil {
		return fmt.Errorf("failed to set RLS context: %w", err)
	}
	return nil
}

// WithRLSContext acquires a connection and sets the RLS context.
// The returned connection must be released by calling conn.Release().
// Usage:
//
//	conn, err := db.WithRLSContext(ctx, pool)
//	if err != nil { return err }
//	defer conn.Release()
//	// Use conn for queries...
func WithRLSContext(ctx context.Context, pool *pgxpool.Pool) (*pgxpool.Conn, error) {
	orgID, ok := GetOrganizationID(ctx)
	if !ok {
		return nil, fmt.Errorf("organization_id not found in context")
	}

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire connection: %w", err)
	}

	if err := SetRLSContext(ctx, conn, orgID); err != nil {
		conn.Release()
		return nil, err
	}

	return conn, nil
}

// RLSPool wraps pgxpool.Pool to automatically apply RLS context from context.Context.
// It implements the DB interface used by repositories.
type RLSPool struct {
	pool *pgxpool.Pool
}

// NewRLSPool creates a new RLS-aware pool wrapper
func NewRLSPool(pool *pgxpool.Pool) *RLSPool {
	return &RLSPool{pool: pool}
}

// Pool returns the underlying pgxpool.Pool
func (r *RLSPool) Pool() *pgxpool.Pool {
	return r.pool
}

// QueryRow executes a query with RLS context and returns a single row.
// If organization_id is in context, it sets the RLS context before querying.
func (r *RLSPool) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	orgID, hasOrg := GetOrganizationID(ctx)
	if !hasOrg {
		// No RLS context, use pool directly
		return r.pool.QueryRow(ctx, sql, args...)
	}

	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return &errorRow{err: fmt.Errorf("failed to acquire connection: %w", err)}
	}

	if err := SetRLSContext(ctx, conn, orgID); err != nil {
		conn.Release()
		return &errorRow{err: err}
	}

	// Return a row that will release the connection after Scan
	return &rlsRow{
		conn: conn,
		row:  conn.QueryRow(ctx, sql, args...),
	}
}

// Query executes a query with RLS context and returns rows.
// If organization_id is in context, it sets the RLS context before querying.
func (r *RLSPool) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	orgID, hasOrg := GetOrganizationID(ctx)
	if !hasOrg {
		// No RLS context, use pool directly
		return r.pool.Query(ctx, sql, args...)
	}

	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire connection: %w", err)
	}

	if err := SetRLSContext(ctx, conn, orgID); err != nil {
		conn.Release()
		return nil, err
	}

	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		conn.Release()
		return nil, err
	}

	// Return rows that will release connection when closed
	return &rlsRows{
		conn: conn,
		Rows: rows,
	}, nil
}

// Exec executes a command with RLS context.
// If organization_id is in context, it sets the RLS context before executing.
func (r *RLSPool) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	orgID, hasOrg := GetOrganizationID(ctx)
	if !hasOrg {
		// No RLS context, use pool directly
		return r.pool.Exec(ctx, sql, args...)
	}

	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	if err := SetRLSContext(ctx, conn, orgID); err != nil {
		return pgconn.CommandTag{}, err
	}

	return conn.Exec(ctx, sql, args...)
}

// rlsRow wraps pgx.Row to release connection after Scan
type rlsRow struct {
	conn *pgxpool.Conn
	row  pgx.Row
}

func (r *rlsRow) Scan(dest ...any) error {
	defer r.conn.Release()
	return r.row.Scan(dest...)
}

// rlsRows wraps pgx.Rows to release connection when closed
type rlsRows struct {
	conn *pgxpool.Conn
	pgx.Rows
}

func (r *rlsRows) Close() {
	r.Rows.Close()
	r.conn.Release()
}

// errorRow implements pgx.Row for error cases
type errorRow struct {
	err error
}

func (r *errorRow) Scan(dest ...any) error {
	return r.err
}

// Tx is the transaction interface
type Tx interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// rlsTx wraps pgx.Tx with connection release on commit/rollback
type rlsTx struct {
	conn *pgxpool.Conn
	tx   pgx.Tx
}

func (t *rlsTx) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return t.tx.QueryRow(ctx, sql, args...)
}

func (t *rlsTx) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return t.tx.Query(ctx, sql, args...)
}

func (t *rlsTx) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return t.tx.Exec(ctx, sql, args...)
}

func (t *rlsTx) Commit(ctx context.Context) error {
	err := t.tx.Commit(ctx)
	t.conn.Release()
	return err
}

func (t *rlsTx) Rollback(ctx context.Context) error {
	err := t.tx.Rollback(ctx)
	t.conn.Release()
	return err
}

// Begin starts a new transaction. RLS context is NOT applied to transactions
// since organization management operations typically need to bypass RLS.
func (r *RLSPool) Begin(ctx context.Context) (Tx, error) {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire connection: %w", err)
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		conn.Release()
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	return &rlsTx{conn: conn, tx: tx}, nil
}
