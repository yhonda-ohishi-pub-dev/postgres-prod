//go:build integration

package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestIntegration_FlickrPhoto_CRUD(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	photoRepo := NewFlickrPhotoRepository(pool)
	orgRepo := NewOrganizationRepository(pool)
	ctx := context.Background()

	// Create a test organization first
	uniqueSlug := fmt.Sprintf("test-flickr-%s", uuid.New().String()[:8])
	org, err := orgRepo.Create(ctx, "Test Org for Photos", uniqueSlug)
	if err != nil {
		t.Fatalf("Failed to create test organization: %v", err)
	}
	defer orgRepo.Delete(ctx, org.ID) // Clean up organization at the end

	// 1. Create
	t.Run("Create", func(t *testing.T) {
		photo, err := photoRepo.Create(ctx, "12345678901", org.ID, "secret123", "server456")
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
		if photo.ID != "12345678901" {
			t.Errorf("Create: ID = %s, want 12345678901", photo.ID)
		}
		if photo.OrganizationID != org.ID {
			t.Errorf("Create: OrganizationID = %s, want %s", photo.OrganizationID, org.ID)
		}
		if photo.Secret != "secret123" {
			t.Errorf("Create: Secret = %s, want secret123", photo.Secret)
		}
		if photo.Server != "server456" {
			t.Errorf("Create: Server = %s, want server456", photo.Server)
		}
		fmt.Printf("✓ Create: ID=%s, OrganizationID=%s, Secret=%s, Server=%s\n", photo.ID, photo.OrganizationID, photo.Secret, photo.Server)

		// 2. GetByID
		t.Run("GetByID", func(t *testing.T) {
			fetched, err := photoRepo.GetByID(ctx, photo.ID)
			if err != nil {
				t.Fatalf("GetByID failed: %v", err)
			}
			if fetched.ID != photo.ID {
				t.Errorf("GetByID: ID = %s, want %s", fetched.ID, photo.ID)
			}
			if fetched.OrganizationID != photo.OrganizationID {
				t.Errorf("GetByID: OrganizationID = %s, want %s", fetched.OrganizationID, photo.OrganizationID)
			}
			if fetched.Secret != photo.Secret {
				t.Errorf("GetByID: Secret = %s, want %s", fetched.Secret, photo.Secret)
			}
			if fetched.Server != photo.Server {
				t.Errorf("GetByID: Server = %s, want %s", fetched.Server, photo.Server)
			}
			fmt.Printf("✓ GetByID: ID=%s, OrganizationID=%s\n", fetched.ID, fetched.OrganizationID)
		})

		// 3. Update
		t.Run("Update", func(t *testing.T) {
			updated, err := photoRepo.Update(ctx, photo.ID, "newsecret999", "newserver888")
			if err != nil {
				t.Fatalf("Update failed: %v", err)
			}
			if updated.Secret != "newsecret999" {
				t.Errorf("Update: Secret = %s, want newsecret999", updated.Secret)
			}
			if updated.Server != "newserver888" {
				t.Errorf("Update: Server = %s, want newserver888", updated.Server)
			}
			if updated.OrganizationID != photo.OrganizationID {
				t.Errorf("Update: OrganizationID should not change, got %s", updated.OrganizationID)
			}
			fmt.Printf("✓ Update: ID=%s, Secret=%s, Server=%s\n", updated.ID, updated.Secret, updated.Server)
		})

		// 4. Create additional photos for List tests
		photo2, err := photoRepo.Create(ctx, "22222222222", org.ID, "secret222", "server222")
		if err != nil {
			t.Fatalf("Create second photo failed: %v", err)
		}
		defer photoRepo.Delete(ctx, photo2.ID)

		photo3, err := photoRepo.Create(ctx, "33333333333", org.ID, "secret333", "server333")
		if err != nil {
			t.Fatalf("Create third photo failed: %v", err)
		}
		defer photoRepo.Delete(ctx, photo3.ID)

		// 5. List
		t.Run("List", func(t *testing.T) {
			photos, err := photoRepo.List(ctx, 10, 0)
			if err != nil {
				t.Fatalf("List failed: %v", err)
			}
			if len(photos) < 3 {
				t.Errorf("List: expected at least 3 photos, got %d", len(photos))
			}
			fmt.Printf("✓ List: returned %d photos\n", len(photos))
		})

		// 6. ListByOrganization
		t.Run("ListByOrganization", func(t *testing.T) {
			photos, err := photoRepo.ListByOrganization(ctx, org.ID, 10, 0)
			if err != nil {
				t.Fatalf("ListByOrganization failed: %v", err)
			}
			if len(photos) != 3 {
				t.Errorf("ListByOrganization: expected 3 photos for org, got %d", len(photos))
			}
			// Verify all photos belong to the organization
			for _, p := range photos {
				if p.OrganizationID != org.ID {
					t.Errorf("ListByOrganization: photo %s has wrong organization_id %s", p.ID, p.OrganizationID)
				}
			}
			fmt.Printf("✓ ListByOrganization: returned %d photos for organization %s\n", len(photos), org.ID)
		})

		// 7. Delete
		t.Run("Delete", func(t *testing.T) {
			err := photoRepo.Delete(ctx, photo.ID)
			if err != nil {
				t.Fatalf("Delete failed: %v", err)
			}

			// Verify hard delete
			_, err = photoRepo.GetByID(ctx, photo.ID)
			if err != ErrFlickrPhotoNotFound {
				t.Errorf("Delete: GetByID should return ErrFlickrPhotoNotFound, got %v", err)
			}
			fmt.Printf("✓ Delete: photo hard deleted\n")
		})
	})
}
