//go:build integration

package repository

import (
	"context"
	"fmt"
	"testing"
)

func TestIntegration_CamFiles_CRUD(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewCamFileRepository(pool)
	orgRepo := NewOrganizationRepository(pool)
	ctx := context.Background()

	// Create a test organization first (slug is auto-generated)
	testOrg, err := orgRepo.Create(ctx, "Test Org for CamFiles")
	if err != nil {
		t.Fatalf("Failed to create test organization: %v", err)
	}
	defer orgRepo.Delete(ctx, testOrg.ID) // Clean up organization at the end

	// Test data
	testName := "test_file_001.jpg"
	testDate := "2025-12-07"
	testHour := "14"
	testType := "image"
	testCam := "cam01"
	testFlickrID := "flickr123456"

	// 1. Create
	t.Run("Create", func(t *testing.T) {
		file, err := repo.Create(ctx, testName, testOrg.ID, testDate, testHour, testType, testCam, &testFlickrID)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
		if file.Name != testName {
			t.Errorf("Create: Name = %s, want %s", file.Name, testName)
		}
		if file.OrganizationID != testOrg.ID {
			t.Errorf("Create: OrganizationID = %s, want %s", file.OrganizationID, testOrg.ID)
		}
		if file.Date != testDate {
			t.Errorf("Create: Date = %s, want %s", file.Date, testDate)
		}
		if file.Hour != testHour {
			t.Errorf("Create: Hour = %s, want %s", file.Hour, testHour)
		}
		if file.Type != testType {
			t.Errorf("Create: Type = %s, want %s", file.Type, testType)
		}
		if file.Cam != testCam {
			t.Errorf("Create: Cam = %s, want %s", file.Cam, testCam)
		}
		if file.FlickrID == nil || *file.FlickrID != testFlickrID {
			t.Errorf("Create: FlickrID = %v, want %s", file.FlickrID, testFlickrID)
		}
		fmt.Printf("✓ Create: Name=%s, OrgID=%s, Date=%s, Hour=%s, Type=%s, Cam=%s, FlickrID=%s\n",
			file.Name, file.OrganizationID, file.Date, file.Hour, file.Type, file.Cam, *file.FlickrID)

		// 2. GetByNameAndOrg
		t.Run("GetByNameAndOrg", func(t *testing.T) {
			fetched, err := repo.GetByNameAndOrg(ctx, testName, testOrg.ID)
			if err != nil {
				t.Fatalf("GetByNameAndOrg failed: %v", err)
			}
			if fetched.Name != file.Name {
				t.Errorf("GetByNameAndOrg: Name = %s, want %s", fetched.Name, file.Name)
			}
			if fetched.OrganizationID != file.OrganizationID {
				t.Errorf("GetByNameAndOrg: OrganizationID = %s, want %s", fetched.OrganizationID, file.OrganizationID)
			}
			if fetched.Date != file.Date {
				t.Errorf("GetByNameAndOrg: Date = %s, want %s", fetched.Date, file.Date)
			}
			fmt.Printf("✓ GetByNameAndOrg: Name=%s, OrgID=%s, Date=%s\n", fetched.Name, fetched.OrganizationID, fetched.Date)
		})

		// 3. Update
		t.Run("Update", func(t *testing.T) {
			updatedDate := "2025-12-08"
			updatedHour := "15"
			updatedType := "video"
			updatedCam := "cam02"
			updatedFlickrID := "flickr789012"

			updated, err := repo.Update(ctx, testName, testOrg.ID, updatedDate, updatedHour, updatedType, updatedCam, &updatedFlickrID)
			if err != nil {
				t.Fatalf("Update failed: %v", err)
			}
			if updated.Name != testName {
				t.Errorf("Update: Name should not change, got %s", updated.Name)
			}
			if updated.OrganizationID != testOrg.ID {
				t.Errorf("Update: OrganizationID should not change, got %s", updated.OrganizationID)
			}
			if updated.Date != updatedDate {
				t.Errorf("Update: Date = %s, want %s", updated.Date, updatedDate)
			}
			if updated.Hour != updatedHour {
				t.Errorf("Update: Hour = %s, want %s", updated.Hour, updatedHour)
			}
			if updated.Type != updatedType {
				t.Errorf("Update: Type = %s, want %s", updated.Type, updatedType)
			}
			if updated.Cam != updatedCam {
				t.Errorf("Update: Cam = %s, want %s", updated.Cam, updatedCam)
			}
			if updated.FlickrID == nil || *updated.FlickrID != updatedFlickrID {
				t.Errorf("Update: FlickrID = %v, want %s", updated.FlickrID, updatedFlickrID)
			}
			fmt.Printf("✓ Update: Name=%s, Date=%s, Hour=%s, Type=%s, Cam=%s, FlickrID=%s\n",
				updated.Name, updated.Date, updated.Hour, updated.Type, updated.Cam, *updated.FlickrID)
		})

		// 4. ListByOrganization
		t.Run("ListByOrganization", func(t *testing.T) {
			// Create a second file for the same organization
			file2Name := "test_file_002.jpg"
			_, err := repo.Create(ctx, file2Name, testOrg.ID, "2025-12-09", "16", "image", "cam03", nil)
			if err != nil {
				t.Fatalf("Failed to create second file: %v", err)
			}

			files, err := repo.ListByOrganization(ctx, testOrg.ID, 10, 0)
			if err != nil {
				t.Fatalf("ListByOrganization failed: %v", err)
			}
			if len(files) < 2 {
				t.Errorf("ListByOrganization: should return at least 2 files, got %d", len(files))
			}
			fmt.Printf("✓ ListByOrganization: returned %d files for organization %s\n", len(files), testOrg.ID)

			// Clean up second file
			repo.Delete(ctx, file2Name, testOrg.ID)
		})

		// 5. Delete
		t.Run("Delete", func(t *testing.T) {
			err := repo.Delete(ctx, testName, testOrg.ID)
			if err != nil {
				t.Fatalf("Delete failed: %v", err)
			}

			// Verify hard delete
			_, err = repo.GetByNameAndOrg(ctx, testName, testOrg.ID)
			if err != ErrCamFileNotFound {
				t.Errorf("Delete: GetByNameAndOrg should return ErrCamFileNotFound, got %v", err)
			}
			fmt.Printf("✓ Delete: cam file hard deleted\n")
		})
	})
}

func TestIntegration_CamFiles_NullFlickrID(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewCamFileRepository(pool)
	orgRepo := NewOrganizationRepository(pool)
	ctx := context.Background()

	// Create a test organization (slug is auto-generated)
	testOrg, err := orgRepo.Create(ctx, "Test Org for Null FlickrID")
	if err != nil {
		t.Fatalf("Failed to create test organization: %v", err)
	}
	defer orgRepo.Delete(ctx, testOrg.ID)

	// Test with null FlickrID
	testName := "test_file_null_flickr.jpg"
	file, err := repo.Create(ctx, testName, testOrg.ID, "2025-12-07", "14", "image", "cam01", nil)
	if err != nil {
		t.Fatalf("Create with null FlickrID failed: %v", err)
	}
	defer repo.Delete(ctx, testName, testOrg.ID)

	if file.FlickrID != nil {
		t.Errorf("FlickrID should be nil, got %v", file.FlickrID)
	}
	fmt.Printf("✓ Create with null FlickrID: Name=%s, FlickrID=nil\n", file.Name)

	// Verify retrieval
	fetched, err := repo.GetByNameAndOrg(ctx, testName, testOrg.ID)
	if err != nil {
		t.Fatalf("GetByNameAndOrg failed: %v", err)
	}
	if fetched.FlickrID != nil {
		t.Errorf("Retrieved FlickrID should be nil, got %v", fetched.FlickrID)
	}
	fmt.Printf("✓ GetByNameAndOrg with null FlickrID verified\n")
}
