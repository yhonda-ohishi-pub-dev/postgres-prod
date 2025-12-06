//go:build integration

package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestIntegration_CamFileExe_CRUD(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewCamFileExeRepository(pool)
	orgRepo := NewOrganizationRepository(pool)
	ctx := context.Background()

	// Create a test organization first
	uniqueSlug := fmt.Sprintf("test-camfileexe-%s", uuid.New().String()[:8])
	org, err := orgRepo.Create(ctx, "Test CamFileExe Org", uniqueSlug)
	if err != nil {
		t.Fatalf("Failed to create test organization: %v", err)
	}
	defer orgRepo.Delete(ctx, org.ID)

	// 1. Create
	t.Run("Create", func(t *testing.T) {
		camFileExe, err := repo.Create(ctx, "test-file.nc", "CAM-001", org.ID, 1)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
		if camFileExe.Name != "test-file.nc" {
			t.Errorf("Create: Name = %s, want test-file.nc", camFileExe.Name)
		}
		if camFileExe.Cam != "CAM-001" {
			t.Errorf("Create: Cam = %s, want CAM-001", camFileExe.Cam)
		}
		if camFileExe.OrganizationID != org.ID {
			t.Errorf("Create: OrganizationID = %s, want %s", camFileExe.OrganizationID, org.ID)
		}
		if camFileExe.Stage != 1 {
			t.Errorf("Create: Stage = %d, want 1", camFileExe.Stage)
		}
		fmt.Printf("✓ Create: Name=%s, Cam=%s, OrganizationID=%s, Stage=%d\n",
			camFileExe.Name, camFileExe.Cam, camFileExe.OrganizationID, camFileExe.Stage)

		// 2. GetByKey
		t.Run("GetByKey", func(t *testing.T) {
			fetched, err := repo.GetByKey(ctx, "test-file.nc", "CAM-001", org.ID)
			if err != nil {
				t.Fatalf("GetByKey failed: %v", err)
			}
			if fetched.Name != camFileExe.Name {
				t.Errorf("GetByKey: Name = %s, want %s", fetched.Name, camFileExe.Name)
			}
			if fetched.Cam != camFileExe.Cam {
				t.Errorf("GetByKey: Cam = %s, want %s", fetched.Cam, camFileExe.Cam)
			}
			if fetched.Stage != camFileExe.Stage {
				t.Errorf("GetByKey: Stage = %d, want %d", fetched.Stage, camFileExe.Stage)
			}
			fmt.Printf("✓ GetByKey: Name=%s, Cam=%s, Stage=%d\n", fetched.Name, fetched.Cam, fetched.Stage)
		})

		// 3. Update
		t.Run("Update", func(t *testing.T) {
			updated, err := repo.Update(ctx, "test-file.nc", "CAM-001", org.ID, 2)
			if err != nil {
				t.Fatalf("Update failed: %v", err)
			}
			if updated.Stage != 2 {
				t.Errorf("Update: Stage = %d, want 2", updated.Stage)
			}
			fmt.Printf("✓ Update: Name=%s, Cam=%s, Stage=%d\n", updated.Name, updated.Cam, updated.Stage)
		})

		// 4. Create more entries for list test
		t.Run("CreateAdditionalEntries", func(t *testing.T) {
			_, err := repo.Create(ctx, "test-file2.nc", "CAM-002", org.ID, 1)
			if err != nil {
				t.Fatalf("Create additional entry 1 failed: %v", err)
			}
			_, err = repo.Create(ctx, "test-file3.nc", "CAM-003", org.ID, 3)
			if err != nil {
				t.Fatalf("Create additional entry 2 failed: %v", err)
			}
			fmt.Printf("✓ Created additional entries for list test\n")
		})

		// 5. ListByOrganization
		t.Run("ListByOrganization", func(t *testing.T) {
			camFileExes, err := repo.ListByOrganization(ctx, org.ID, 10, 0)
			if err != nil {
				t.Fatalf("ListByOrganization failed: %v", err)
			}
			if len(camFileExes) < 3 {
				t.Errorf("ListByOrganization: should return at least 3 entries, got %d", len(camFileExes))
			}
			fmt.Printf("✓ ListByOrganization: returned %d entries\n", len(camFileExes))
		})

		// 6. Delete
		t.Run("Delete", func(t *testing.T) {
			err := repo.Delete(ctx, "test-file.nc", "CAM-001", org.ID)
			if err != nil {
				t.Fatalf("Delete failed: %v", err)
			}

			// Verify deletion
			_, err = repo.GetByKey(ctx, "test-file.nc", "CAM-001", org.ID)
			if err != ErrCamFileExeNotFound {
				t.Errorf("Delete: GetByKey should return ErrCamFileExeNotFound, got %v", err)
			}
			fmt.Printf("✓ Delete: cam file exe entry deleted\n")
		})

		// 7. GetByKey on non-existent entry
		t.Run("GetByKey_NotFound", func(t *testing.T) {
			_, err := repo.GetByKey(ctx, "non-existent.nc", "CAM-999", org.ID)
			if err != ErrCamFileExeNotFound {
				t.Errorf("GetByKey_NotFound: should return ErrCamFileExeNotFound, got %v", err)
			}
			fmt.Printf("✓ GetByKey_NotFound: correctly returned ErrCamFileExeNotFound\n")
		})

		// 8. Update on non-existent entry
		t.Run("Update_NotFound", func(t *testing.T) {
			_, err := repo.Update(ctx, "non-existent.nc", "CAM-999", org.ID, 5)
			if err != ErrCamFileExeNotFound {
				t.Errorf("Update_NotFound: should return ErrCamFileExeNotFound, got %v", err)
			}
			fmt.Printf("✓ Update_NotFound: correctly returned ErrCamFileExeNotFound\n")
		})

		// 9. Delete on non-existent entry
		t.Run("Delete_NotFound", func(t *testing.T) {
			err := repo.Delete(ctx, "non-existent.nc", "CAM-999", org.ID)
			if err != ErrCamFileExeNotFound {
				t.Errorf("Delete_NotFound: should return ErrCamFileExeNotFound, got %v", err)
			}
			fmt.Printf("✓ Delete_NotFound: correctly returned ErrCamFileExeNotFound\n")
		})

		// Cleanup remaining test entries
		repo.Delete(ctx, "test-file2.nc", "CAM-002", org.ID)
		repo.Delete(ctx, "test-file3.nc", "CAM-003", org.ID)
	})
}

func TestIntegration_CamFileExe_ListPagination(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewCamFileExeRepository(pool)
	orgRepo := NewOrganizationRepository(pool)
	ctx := context.Background()

	// Create a test organization
	paginationSlug := fmt.Sprintf("test-pagination-%s", uuid.New().String()[:8])
	org, err := orgRepo.Create(ctx, "Test Pagination Org", paginationSlug)
	if err != nil {
		t.Fatalf("Failed to create test organization: %v", err)
	}
	defer orgRepo.Delete(ctx, org.ID)

	// Create multiple entries for pagination test
	for i := 1; i <= 15; i++ {
		_, err := repo.Create(ctx, fmt.Sprintf("file-%02d.nc", i), fmt.Sprintf("CAM-%03d", i), org.ID, int32(i))
		if err != nil {
			t.Fatalf("Failed to create test entry %d: %v", i, err)
		}
	}
	defer func() {
		for i := 1; i <= 15; i++ {
			repo.Delete(ctx, fmt.Sprintf("file-%02d.nc", i), fmt.Sprintf("CAM-%03d", i), org.ID)
		}
	}()

	t.Run("Pagination_FirstPage", func(t *testing.T) {
		camFileExes, err := repo.ListByOrganization(ctx, org.ID, 10, 0)
		if err != nil {
			t.Fatalf("ListByOrganization failed: %v", err)
		}
		if len(camFileExes) != 10 {
			t.Errorf("Pagination_FirstPage: expected 10 entries, got %d", len(camFileExes))
		}
		fmt.Printf("✓ Pagination_FirstPage: returned %d entries\n", len(camFileExes))
	})

	t.Run("Pagination_SecondPage", func(t *testing.T) {
		camFileExes, err := repo.ListByOrganization(ctx, org.ID, 10, 10)
		if err != nil {
			t.Fatalf("ListByOrganization failed: %v", err)
		}
		if len(camFileExes) != 5 {
			t.Errorf("Pagination_SecondPage: expected 5 entries, got %d", len(camFileExes))
		}
		fmt.Printf("✓ Pagination_SecondPage: returned %d entries\n", len(camFileExes))
	})

	t.Run("Pagination_DefaultLimit", func(t *testing.T) {
		camFileExes, err := repo.ListByOrganization(ctx, org.ID, 0, 0)
		if err != nil {
			t.Fatalf("ListByOrganization failed: %v", err)
		}
		if len(camFileExes) != 10 {
			t.Errorf("Pagination_DefaultLimit: expected 10 entries (default), got %d", len(camFileExes))
		}
		fmt.Printf("✓ Pagination_DefaultLimit: returned %d entries\n", len(camFileExes))
	})

	t.Run("Pagination_MaxLimit", func(t *testing.T) {
		camFileExes, err := repo.ListByOrganization(ctx, org.ID, 200, 0)
		if err != nil {
			t.Fatalf("ListByOrganization failed: %v", err)
		}
		if len(camFileExes) > 100 {
			t.Errorf("Pagination_MaxLimit: should enforce max limit of 100, got %d", len(camFileExes))
		}
		fmt.Printf("✓ Pagination_MaxLimit: correctly enforced max limit, returned %d entries\n", len(camFileExes))
	})
}
