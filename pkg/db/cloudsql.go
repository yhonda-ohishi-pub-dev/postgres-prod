package db

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewCloudSQLPool creates a connection pool to Cloud SQL using IAM authentication.
// This uses the Cloud SQL Go Connector which automatically handles:
// - IAM database authentication
// - Secure TLS connections
// - Automatic credential refresh
func NewCloudSQLPool(ctx context.Context, instanceConnection, dbUser, dbName string) (*pgxpool.Pool, func() error, error) {
	// Create Cloud SQL dialer with IAM authentication
	dialer, err := cloudsqlconn.NewDialer(ctx, cloudsqlconn.WithIAMAuthN())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create Cloud SQL dialer: %w", err)
	}

	// Configure connection pool
	dsn := fmt.Sprintf("user=%s database=%s", dbUser, dbName)
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		dialer.Close()
		return nil, nil, fmt.Errorf("failed to parse pgx config: %w", err)
	}

	// Use Cloud SQL dialer for connections
	config.ConnConfig.DialFunc = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dialer.Dial(ctx, instanceConnection)
	}

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		dialer.Close()
		return nil, nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		dialer.Close()
		return nil, nil, fmt.Errorf("failed to ping database: %w", err)
	}

	cleanup := func() error {
		pool.Close()
		return dialer.Close()
	}

	return pool, cleanup, nil
}

// NewLocalPool creates a connection pool via Cloud SQL Proxy (for local development).
// Connects to localhost where Cloud SQL Proxy is listening.
func NewLocalPool(ctx context.Context, dbUser, dbName, dbPassword string, port int) (*pgxpool.Pool, error) {
	if port == 0 {
		port = 5432
	}

	// Use URL format to properly handle special characters in username (e.g., @)
	dsn := fmt.Sprintf("postgres://%s:%s@localhost:%d/%s?sslmode=disable",
		url.QueryEscape(dbUser), url.QueryEscape(dbPassword), port, dbName)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}

// NewPool creates a connection pool based on environment.
// Uses Cloud SQL IAM auth in production (K_SERVICE set), local proxy otherwise.
func NewPool(ctx context.Context, instanceConnection, dbUser, dbName, dbPassword string, port int) (*pgxpool.Pool, func() error, error) {
	// Check if running on Cloud Run
	if os.Getenv("K_SERVICE") != "" {
		return NewCloudSQLPool(ctx, instanceConnection, dbUser, dbName)
	}

	// Local development via proxy with --auto-iam-authn (password can be empty)
	if port == 0 {
		port = 5432
	}
	pool, err := NewLocalPool(ctx, dbUser, dbName, dbPassword, port)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() error {
		pool.Close()
		return nil
	}

	return pool, cleanup, nil
}
