//go:build integration

package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestIntegration_Dtakologs_CRUD(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewDtakologsRepository(pool)
	ctx := context.Background()

	// Create unique test data
	orgID := fmt.Sprintf("test-org-%s", uuid.New().String()[:8])
	dataDateTime := time.Now().Format(time.RFC3339)
	var vehicleCd int32 = 12345
	testData := &Dtakologs{
		OrganizationID:         orgID,
		Type:                   "test-type",
		DataDateTime:           dataDateTime,
		VehicleCd:              vehicleCd,
		VehicleName:            "Test Vehicle",
		BranchName:             "Test Branch",
		AllStateRyoutColor:     "green",
		SettingTemp:            "20",
		SettingTemp1:           "21",
		SettingTemp3:           "22",
		SettingTemp4:           "23",
		StateFlag:              "1",
	}

	// 1. Create
	t.Run("Create", func(t *testing.T) {
		err := repo.Create(ctx, testData)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
		fmt.Printf("✓ Create: OrgID=%s, DataDateTime=%s, VehicleCd=%d\n", orgID, dataDateTime, vehicleCd)

		// 2. GetByPrimaryKey
		t.Run("GetByPrimaryKey", func(t *testing.T) {
			fetched, err := repo.GetByPrimaryKey(ctx, orgID, dataDateTime, vehicleCd)
			if err != nil {
				t.Fatalf("GetByPrimaryKey failed: %v", err)
			}
			if fetched.OrganizationID != orgID {
				t.Errorf("GetByPrimaryKey: OrganizationID = %s, want %s", fetched.OrganizationID, orgID)
			}
			if fetched.VehicleCd != vehicleCd {
				t.Errorf("GetByPrimaryKey: VehicleCd = %d, want %d", fetched.VehicleCd, vehicleCd)
			}
			fmt.Printf("✓ GetByPrimaryKey: OrgID=%s, VehicleName=%s\n", fetched.OrganizationID, fetched.VehicleName)
		})

		// 3. Update
		t.Run("Update", func(t *testing.T) {
			testData.VehicleName = "Updated Vehicle"
			err := repo.Update(ctx, testData)
			if err != nil {
				t.Fatalf("Update failed: %v", err)
			}
			// Verify update
			updated, err := repo.GetByPrimaryKey(ctx, orgID, dataDateTime, vehicleCd)
			if err != nil {
				t.Fatalf("GetByPrimaryKey after Update failed: %v", err)
			}
			if updated.VehicleName != "Updated Vehicle" {
				t.Errorf("Update: VehicleName = %s, want %s", updated.VehicleName, "Updated Vehicle")
			}
			fmt.Printf("✓ Update: VehicleName=%s\n", updated.VehicleName)
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
			err := repo.Delete(ctx, orgID, dataDateTime, vehicleCd)
			if err != nil {
				t.Fatalf("Delete failed: %v", err)
			}

			// Verify deletion
			_, err = repo.GetByPrimaryKey(ctx, orgID, dataDateTime, vehicleCd)
			if err != ErrDtakologsNotFound {
				t.Errorf("Delete: GetByPrimaryKey should return ErrDtakologsNotFound, got %v", err)
			}
			fmt.Printf("✓ Delete: record deleted\n")
		})
	})
}
