//go:build integration

package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Test organization IDs - these must exist in the database
// Using pre-seeded organizations from test fixtures
const (
	testOrgA = "11111111-1111-1111-1111-111111111111" // ACME Corp
	testOrgB = "22222222-2222-2222-2222-222222222222" // Globex Inc
)

func setupTestPool(t *testing.T) *pgxpool.Pool {
	ctx := context.Background()
	dsn := "postgres://postgres:postgres@localhost:5432/myapp_postgres?sslmode=disable"
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	return pool
}

// setupRLSTestPool creates a connection pool using a non-superuser account
// so that RLS policies are enforced
func setupRLSTestPool(t *testing.T) *pgxpool.Pool {
	ctx := context.Background()
	// Use app_user which is a non-superuser, so RLS is enforced
	dsn := "postgres://app_user:app_password@localhost:5432/myapp_postgres?sslmode=disable"
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("Failed to connect to database as app_user: %v", err)
	}
	return pool
}

// TestRLS_IsolationBetweenOrganizations verifies that RLS policies correctly
// isolate data between different organizations.
func TestRLS_IsolationBetweenOrganizations(t *testing.T) {
	// Use superuser pool for setup/cleanup (bypasses RLS)
	superPool := setupTestPool(t)
	defer superPool.Close()

	// Use app_user pool for RLS tests (RLS is enforced)
	appPool := setupRLSTestPool(t)
	defer appPool.Close()

	ctx := context.Background()

	// Use existing test organizations
	orgA := testOrgA
	orgB := testOrgB

	// Insert test data for both organizations using superadmin (bypasses RLS)
	testFileA := fmt.Sprintf("test-file-org-a-%s", uuid.New().String()[:8])
	testFileB := fmt.Sprintf("test-file-org-b-%s", uuid.New().String()[:8])

	_, err := superPool.Exec(ctx, `
		INSERT INTO cam_files (name, organization_id, date, hour, type, cam, flickr_id)
		VALUES ($1, $2, '2024-01-01', 12, 'image', 'cam1', NULL)
	`, testFileA, orgA)
	if err != nil {
		t.Fatalf("Failed to insert test data for org A: %v", err)
	}

	_, err = superPool.Exec(ctx, `
		INSERT INTO cam_files (name, organization_id, date, hour, type, cam, flickr_id)
		VALUES ($1, $2, '2024-01-01', 12, 'image', 'cam1', NULL)
	`, testFileB, orgB)
	if err != nil {
		t.Fatalf("Failed to insert test data for org B: %v", err)
	}

	// Cleanup at the end using superuser
	defer func() {
		superPool.Exec(ctx, "DELETE FROM cam_files WHERE name = $1", testFileA)
		superPool.Exec(ctx, "DELETE FROM cam_files WHERE name = $1", testFileB)
	}()

	// Test 1: With org A context, should only see org A's data
	t.Run("OrgA_CanOnlySeeOwnData", func(t *testing.T) {
		rlsPool := NewRLSPool(appPool)
		ctxWithOrgA := WithOrganizationID(ctx, orgA)

		// Should find org A's file
		var count int
		err := rlsPool.QueryRow(ctxWithOrgA,
			"SELECT COUNT(*) FROM cam_files WHERE name = $1", testFileA).Scan(&count)
		if err != nil {
			t.Fatalf("Query failed: %v", err)
		}
		if count != 1 {
			t.Errorf("Expected to find 1 file for org A, got %d", count)
		}

		// Should NOT find org B's file (RLS filters it out)
		err = rlsPool.QueryRow(ctxWithOrgA,
			"SELECT COUNT(*) FROM cam_files WHERE name = $1", testFileB).Scan(&count)
		if err != nil {
			t.Fatalf("Query failed: %v", err)
		}
		if count != 0 {
			t.Errorf("Expected 0 files for org B when querying as org A, got %d", count)
		}

		fmt.Printf("✓ Org A can see own data (%s) but cannot see Org B's data (%s)\n", testFileA, testFileB)
	})

	// Test 2: With org B context, should only see org B's data
	t.Run("OrgB_CanOnlySeeOwnData", func(t *testing.T) {
		rlsPool := NewRLSPool(appPool)
		ctxWithOrgB := WithOrganizationID(ctx, orgB)

		// Should find org B's file
		var count int
		err := rlsPool.QueryRow(ctxWithOrgB,
			"SELECT COUNT(*) FROM cam_files WHERE name = $1", testFileB).Scan(&count)
		if err != nil {
			t.Fatalf("Query failed: %v", err)
		}
		if count != 1 {
			t.Errorf("Expected to find 1 file for org B, got %d", count)
		}

		// Should NOT find org A's file (RLS filters it out)
		err = rlsPool.QueryRow(ctxWithOrgB,
			"SELECT COUNT(*) FROM cam_files WHERE name = $1", testFileA).Scan(&count)
		if err != nil {
			t.Fatalf("Query failed: %v", err)
		}
		if count != 0 {
			t.Errorf("Expected 0 files for org A when querying as org B, got %d", count)
		}

		fmt.Printf("✓ Org B can see own data (%s) but cannot see Org A's data (%s)\n", testFileB, testFileA)
	})

	// Test 3: List query with RLS should only return matching organization's data
	t.Run("ListQuery_ReturnsOnlyOwnData", func(t *testing.T) {
		rlsPool := NewRLSPool(appPool)
		ctxWithOrgA := WithOrganizationID(ctx, orgA)

		rows, err := rlsPool.Query(ctxWithOrgA,
			"SELECT name, organization_id FROM cam_files WHERE name IN ($1, $2)",
			testFileA, testFileB)
		if err != nil {
			t.Fatalf("Query failed: %v", err)
		}
		defer rows.Close()

		var names []string
		for rows.Next() {
			var name, orgID string
			if err := rows.Scan(&name, &orgID); err != nil {
				t.Fatalf("Scan failed: %v", err)
			}
			names = append(names, name)
			if orgID != orgA {
				t.Errorf("Unexpected organization_id: got %s, want %s", orgID, orgA)
			}
		}

		if len(names) != 1 {
			t.Errorf("Expected 1 result, got %d: %v", len(names), names)
		}
		if len(names) > 0 && names[0] != testFileA {
			t.Errorf("Expected %s, got %s", testFileA, names[0])
		}

		fmt.Printf("✓ List query returns only org A's data: %v\n", names)
	})
}

// TestRLS_ContextPropagation verifies that organization_id context is correctly
// propagated through the RLSPool wrapper.
func TestRLS_ContextPropagation(t *testing.T) {
	pool := setupTestPool(t)
	defer pool.Close()

	ctx := context.Background()
	testOrgID := uuid.New().String()

	t.Run("ContextWithOrganizationID", func(t *testing.T) {
		ctxWithOrg := WithOrganizationID(ctx, testOrgID)
		retrievedOrgID, ok := GetOrganizationID(ctxWithOrg)

		if !ok {
			t.Error("Expected to find organization_id in context")
		}
		if retrievedOrgID != testOrgID {
			t.Errorf("Expected %s, got %s", testOrgID, retrievedOrgID)
		}

		fmt.Printf("✓ Context correctly stores organization_id: %s\n", testOrgID)
	})

	t.Run("ContextWithoutOrganizationID", func(t *testing.T) {
		_, ok := GetOrganizationID(ctx)
		if ok {
			t.Error("Expected no organization_id in plain context")
		}

		fmt.Println("✓ Plain context correctly has no organization_id")
	})

	t.Run("RLSPool_RequiresOrganizationID", func(t *testing.T) {
		_, err := WithRLSContext(ctx, pool)
		if err == nil {
			t.Error("Expected error when calling WithRLSContext without organization_id")
		}

		fmt.Println("✓ WithRLSContext correctly requires organization_id")
	})

	t.Run("RLSPool_AcquiresConnectionWithContext", func(t *testing.T) {
		ctxWithOrg := WithOrganizationID(ctx, testOrgID)
		conn, err := WithRLSContext(ctxWithOrg, pool)
		if err != nil {
			t.Fatalf("WithRLSContext failed: %v", err)
		}
		defer conn.Release()

		// Verify the session variable is set
		var currentOrgID string
		err = conn.QueryRow(ctxWithOrg,
			"SELECT current_setting('app.current_organization_id', true)").Scan(&currentOrgID)
		if err != nil {
			t.Fatalf("Query failed: %v", err)
		}
		if currentOrgID != testOrgID {
			t.Errorf("Session variable mismatch: got %s, want %s", currentOrgID, testOrgID)
		}

		fmt.Printf("✓ Session variable correctly set to: %s\n", currentOrgID)
	})
}

// TestRLS_UpdateDeleteIsolation verifies that UPDATE and DELETE operations
// are also properly isolated by RLS.
func TestRLS_UpdateDeleteIsolation(t *testing.T) {
	// Use superuser pool for setup/cleanup (bypasses RLS)
	superPool := setupTestPool(t)
	defer superPool.Close()

	// Use app_user pool for RLS tests (RLS is enforced)
	appPool := setupRLSTestPool(t)
	defer appPool.Close()

	ctx := context.Background()

	orgA := testOrgA
	orgB := testOrgB

	testFileA := fmt.Sprintf("test-update-a-%s", uuid.New().String()[:8])
	testFileB := fmt.Sprintf("test-update-b-%s", uuid.New().String()[:8])

	// Insert test data using superuser
	_, err := superPool.Exec(ctx, `
		INSERT INTO cam_files (name, organization_id, date, hour, type, cam, flickr_id)
		VALUES ($1, $2, '2024-01-01', '10', 'video', 'cam2', NULL)
	`, testFileA, orgA)
	if err != nil {
		t.Fatalf("Failed to insert test data A: %v", err)
	}

	_, err = superPool.Exec(ctx, `
		INSERT INTO cam_files (name, organization_id, date, hour, type, cam, flickr_id)
		VALUES ($1, $2, '2024-01-01', '10', 'video', 'cam2', NULL)
	`, testFileB, orgB)
	if err != nil {
		t.Fatalf("Failed to insert test data B: %v", err)
	}

	defer func() {
		superPool.Exec(ctx, "DELETE FROM cam_files WHERE name = $1", testFileA)
		superPool.Exec(ctx, "DELETE FROM cam_files WHERE name = $1", testFileB)
	}()

	t.Run("Update_OnlyAffectsOwnData", func(t *testing.T) {
		rlsPool := NewRLSPool(appPool)
		ctxWithOrgA := WithOrganizationID(ctx, orgA)

		// Try to update org B's file from org A's context
		result, err := rlsPool.Exec(ctxWithOrgA,
			"UPDATE cam_files SET hour = '99' WHERE name = $1", testFileB)
		if err != nil {
			t.Fatalf("Exec failed: %v", err)
		}

		// Should affect 0 rows (RLS blocks the update)
		if result.RowsAffected() != 0 {
			t.Errorf("Expected 0 rows affected, got %d", result.RowsAffected())
		}

		// Verify org B's data is unchanged using superuser
		var hour string
		err = superPool.QueryRow(ctx, "SELECT hour FROM cam_files WHERE name = $1", testFileB).Scan(&hour)
		if err != nil {
			t.Fatalf("Verification query failed: %v", err)
		}
		if hour != "10" {
			t.Errorf("Org B's data was modified: hour = %s, want 10", hour)
		}

		fmt.Println("✓ Update from Org A cannot modify Org B's data")
	})

	t.Run("Delete_OnlyAffectsOwnData", func(t *testing.T) {
		rlsPool := NewRLSPool(appPool)
		ctxWithOrgA := WithOrganizationID(ctx, orgA)

		// Try to delete org B's file from org A's context
		result, err := rlsPool.Exec(ctxWithOrgA,
			"DELETE FROM cam_files WHERE name = $1", testFileB)
		if err != nil {
			t.Fatalf("Exec failed: %v", err)
		}

		// Should affect 0 rows (RLS blocks the delete)
		if result.RowsAffected() != 0 {
			t.Errorf("Expected 0 rows affected, got %d", result.RowsAffected())
		}

		// Verify org B's data still exists using superuser
		var count int
		err = superPool.QueryRow(ctx, "SELECT COUNT(*) FROM cam_files WHERE name = $1", testFileB).Scan(&count)
		if err != nil {
			t.Fatalf("Verification query failed: %v", err)
		}
		if count != 1 {
			t.Errorf("Org B's data was deleted: count = %d, want 1", count)
		}

		fmt.Println("✓ Delete from Org A cannot remove Org B's data")
	})
}
