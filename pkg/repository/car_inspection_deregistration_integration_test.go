//go:build integration

package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestIntegration_CarInspectionDeregistration_CRUD(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewCarInspectionDeregistrationRepository(pool)
	ctx := context.Background()

	// Create unique test data
	orgID := fmt.Sprintf("test-org-%s", uuid.New().String()[:8])
	carID := fmt.Sprintf("car-%s", uuid.New().String()[:8])
	twodimensionCodeInfoCarNo := "1234567890"
	carNo := "ABC-1234"
	validPeriodExpirDateE := "2025"
	validPeriodExpirDateY := "2025"
	validPeriodExpirDateM := "12"
	validPeriodExpirDateD := "31"
	twodimensionCodeInfoValidPeriodExpirDate := "2025-12-31"

	// 1. Create
	t.Run("Create", func(t *testing.T) {
		record, err := repo.Create(ctx, orgID, carID, twodimensionCodeInfoCarNo, carNo, validPeriodExpirDateE, validPeriodExpirDateY, validPeriodExpirDateM, validPeriodExpirDateD, twodimensionCodeInfoValidPeriodExpirDate)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
		if record.OrganizationID != orgID {
			t.Errorf("Create: OrganizationID = %s, want %s", record.OrganizationID, orgID)
		}
		if record.CarID != carID {
			t.Errorf("Create: CarID = %s, want %s", record.CarID, carID)
		}
		fmt.Printf("✓ Create: OrgID=%s, CarID=%s\n", record.OrganizationID, record.CarID)

		// 2. GetByPrimaryKey
		t.Run("GetByPrimaryKey", func(t *testing.T) {
			fetched, err := repo.GetByPrimaryKey(ctx, orgID, carID, twodimensionCodeInfoValidPeriodExpirDate)
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

		// 3. Update
		t.Run("Update", func(t *testing.T) {
			updatedCarNo := "XYZ-5678"
			updated, err := repo.Update(ctx, orgID, carID, twodimensionCodeInfoCarNo, updatedCarNo, validPeriodExpirDateE, validPeriodExpirDateY, validPeriodExpirDateM, validPeriodExpirDateD, twodimensionCodeInfoValidPeriodExpirDate)
			if err != nil {
				t.Fatalf("Update failed: %v", err)
			}
			if updated.CarNo != updatedCarNo {
				t.Errorf("Update: CarNo = %s, want %s", updated.CarNo, updatedCarNo)
			}
			fmt.Printf("✓ Update: CarNo=%s\n", updated.CarNo)
		})

		// 4. List
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

		// 5. Delete
		t.Run("Delete", func(t *testing.T) {
			err := repo.Delete(ctx, orgID, carID, twodimensionCodeInfoValidPeriodExpirDate)
			if err != nil {
				t.Fatalf("Delete failed: %v", err)
			}

			// Verify deletion
			_, err = repo.GetByPrimaryKey(ctx, orgID, carID, twodimensionCodeInfoValidPeriodExpirDate)
			if err != ErrCarInspectionDeregistrationNotFound {
				t.Errorf("Delete: GetByPrimaryKey should return ErrCarInspectionDeregistrationNotFound, got %v", err)
			}
			fmt.Printf("✓ Delete: record deleted\n")
		})
	})
}
