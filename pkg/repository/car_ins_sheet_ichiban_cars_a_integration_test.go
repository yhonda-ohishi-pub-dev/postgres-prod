//go:build integration

package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestIntegration_CarInsSheetIchibanCarsA_CRUD(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewCarInsSheetIchibanCarsARepository(pool)
	ctx := context.Background()

	// Create unique test data
	orgID := fmt.Sprintf("test-org-%s", uuid.New().String()[:8])
	electCertMgNo := fmt.Sprintf("cert-%s", uuid.New().String()[:8])
	grantdateE := "R"
	grantdateY := "05"
	grantdateM := "12"
	grantdateD := "07"
	idCars := fmt.Sprintf("car-%s", uuid.New().String()[:8])

	// 1. Create
	t.Run("Create", func(t *testing.T) {
		record, err := repo.Create(ctx, orgID, electCertMgNo, grantdateE, grantdateY, grantdateM, grantdateD, &idCars)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
		if record.OrganizationID != orgID {
			t.Errorf("Create: OrganizationID = %s, want %s", record.OrganizationID, orgID)
		}
		fmt.Printf("✓ Create: OrgID=%s, ElectCertMgNo=%s\n", record.OrganizationID, record.ElectCertMgNo)

		// 2. GetByPrimaryKey
		t.Run("GetByPrimaryKey", func(t *testing.T) {
			fetched, err := repo.GetByPrimaryKey(ctx, orgID, electCertMgNo, grantdateE, grantdateY, grantdateM, grantdateD)
			if err != nil {
				t.Fatalf("GetByPrimaryKey failed: %v", err)
			}
			if fetched.OrganizationID != orgID {
				t.Errorf("GetByPrimaryKey: OrganizationID = %s, want %s", fetched.OrganizationID, orgID)
			}
			fmt.Printf("✓ GetByPrimaryKey: OrgID=%s, ElectCertMgNo=%s\n", fetched.OrganizationID, fetched.ElectCertMgNo)
		})

		// 3. Update
		t.Run("Update", func(t *testing.T) {
			newIdCars := fmt.Sprintf("updated-car-%s", uuid.New().String()[:8])
			updated, err := repo.Update(ctx, orgID, electCertMgNo, grantdateE, grantdateY, grantdateM, grantdateD, &newIdCars)
			if err != nil {
				t.Fatalf("Update failed: %v", err)
			}
			if *updated.IDCars != newIdCars {
				t.Errorf("Update: IDCars = %s, want %s", *updated.IDCars, newIdCars)
			}
			fmt.Printf("✓ Update: IDCars=%s\n", *updated.IDCars)
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
			err := repo.Delete(ctx, orgID, electCertMgNo, grantdateE, grantdateY, grantdateM, grantdateD)
			if err != nil {
				t.Fatalf("Delete failed: %v", err)
			}

			// Verify deletion
			_, err = repo.GetByPrimaryKey(ctx, orgID, electCertMgNo, grantdateE, grantdateY, grantdateM, grantdateD)
			if err != ErrCarInsSheetIchibanCarsANotFound {
				t.Errorf("Delete: GetByPrimaryKey should return ErrCarInsSheetIchibanCarsANotFound, got %v", err)
			}
			fmt.Printf("✓ Delete: record deleted\n")
		})
	})
}
