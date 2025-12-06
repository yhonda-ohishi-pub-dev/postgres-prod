//go:build integration

package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestIntegration_CarInspectionDeregistrationFiles_CRUD(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewCarInspectionDeregistrationFilesRepository(pool)
	ctx := context.Background()

	// Create unique test data
	orgID := fmt.Sprintf("test-org-%s", uuid.New().String()[:8])
	carID := fmt.Sprintf("car-%s", uuid.New().String()[:8])
	expirDate := "2025-12-31"
	fileUUID := uuid.New().String()

	// 1. Create
	t.Run("Create", func(t *testing.T) {
		record, err := repo.Create(ctx, orgID, carID, expirDate, fileUUID)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
		if record.OrganizationID != orgID {
			t.Errorf("Create: OrganizationID = %s, want %s", record.OrganizationID, orgID)
		}
		if record.CarID != carID {
			t.Errorf("Create: CarID = %s, want %s", record.CarID, carID)
		}
		if record.FileUUID != fileUUID {
			t.Errorf("Create: FileUUID = %s, want %s", record.FileUUID, fileUUID)
		}
		fmt.Printf("✓ Create: OrgID=%s, CarID=%s, FileUUID=%s\n", record.OrganizationID, record.CarID, record.FileUUID)

		// 2. GetByPrimaryKey
		t.Run("GetByPrimaryKey", func(t *testing.T) {
			fetched, err := repo.GetByPrimaryKey(ctx, orgID, carID, expirDate, fileUUID)
			if err != nil {
				t.Fatalf("GetByPrimaryKey failed: %v", err)
			}
			if fetched.OrganizationID != orgID {
				t.Errorf("GetByPrimaryKey: OrganizationID = %s, want %s", fetched.OrganizationID, orgID)
			}
			if fetched.CarID != carID {
				t.Errorf("GetByPrimaryKey: CarID = %s, want %s", fetched.CarID, carID)
			}
			fmt.Printf("✓ GetByPrimaryKey: OrgID=%s, CarID=%s\n", fetched.OrganizationID, fetched.CarID)
		})

		// 3. ListByOrganization
		t.Run("ListByOrganization", func(t *testing.T) {
			records, err := repo.ListByOrganization(ctx, orgID, 10, 0)
			if err != nil {
				t.Fatalf("ListByOrganization failed: %v", err)
			}
			if len(records) == 0 {
				t.Error("ListByOrganization: should return at least one record")
			}
			fmt.Printf("✓ ListByOrganization: returned %d records\n", len(records))
		})

		// 4. ListByCarInspectionDeregistration
		t.Run("ListByCarInspectionDeregistration", func(t *testing.T) {
			records, err := repo.ListByCarInspectionDeregistration(ctx, orgID, carID, expirDate, 10, 0)
			if err != nil {
				t.Fatalf("ListByCarInspectionDeregistration failed: %v", err)
			}
			if len(records) == 0 {
				t.Error("ListByCarInspectionDeregistration: should return at least one record")
			}
			fmt.Printf("✓ ListByCarInspectionDeregistration: returned %d records\n", len(records))
		})

		// 5. Delete
		t.Run("Delete", func(t *testing.T) {
			err := repo.Delete(ctx, orgID, carID, expirDate, fileUUID)
			if err != nil {
				t.Fatalf("Delete failed: %v", err)
			}

			// Verify deletion
			_, err = repo.GetByPrimaryKey(ctx, orgID, carID, expirDate, fileUUID)
			if err != ErrCarInspectionDeregistrationFilesNotFound {
				t.Errorf("Delete: GetByPrimaryKey should return ErrCarInspectionDeregistrationFilesNotFound, got %v", err)
			}
			fmt.Printf("✓ Delete: record deleted\n")
		})
	})
}
