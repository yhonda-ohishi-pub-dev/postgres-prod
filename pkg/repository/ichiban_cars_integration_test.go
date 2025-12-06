//go:build integration

package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestIntegration_IchibanCars_CRUD(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewIchibanCarRepository(pool)
	orgRepo := NewOrganizationRepository(pool)
	ctx := context.Background()

	// Setup: Create test organization
	uniqueSlug := fmt.Sprintf("test-ichiban-%s", uuid.New().String()[:8])
	testOrg, err := orgRepo.Create(ctx, "Test IchibanCars Org", uniqueSlug)
	if err != nil {
		t.Fatalf("Setup: failed to create test organization: %v", err)
	}
	defer orgRepo.Delete(ctx, testOrg.ID)

	// 1. Create
	t.Run("Create", func(t *testing.T) {
		id := "CAR001"
		id4 := "1234"
		shashu := "Toyota"
		name := "Test Car"
		nameR := "テストカー"
		sekisai := 5.0
		regDate := "2024-01-01"
		parchDate := "2024-01-15"
		bumonCodeID := "DEPT01"
		driverID := "DRV001"

		car, err := repo.Create(ctx, id, testOrg.ID, id4, shashu, &name, &nameR, &sekisai, &regDate, &parchDate, nil, &bumonCodeID, &driverID)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
		if car.ID != id {
			t.Errorf("Create: ID = %s, want %s", car.ID, id)
		}
		if car.OrganizationID != testOrg.ID {
			t.Errorf("Create: OrganizationID = %s, want %s", car.OrganizationID, testOrg.ID)
		}
		if car.ID4 != id4 {
			t.Errorf("Create: ID4 = %s, want %s", car.ID4, id4)
		}
		if car.Shashu != shashu {
			t.Errorf("Create: Shashu = %s, want %s", car.Shashu, shashu)
		}
		if car.Name == nil || *car.Name != name {
			t.Errorf("Create: Name = %v, want %s", car.Name, name)
		}
		if car.Sekisai == nil || *car.Sekisai != sekisai {
			t.Errorf("Create: Sekisai = %v, want %f", car.Sekisai, sekisai)
		}
		fmt.Printf("✓ Create: ID=%s, OrganizationID=%s, ID4=%s, Shashu=%s\n", car.ID, car.OrganizationID, car.ID4, car.Shashu)

		// 2. GetByIDAndOrg
		t.Run("GetByIDAndOrg", func(t *testing.T) {
			fetched, err := repo.GetByIDAndOrg(ctx, id, testOrg.ID)
			if err != nil {
				t.Fatalf("GetByIDAndOrg failed: %v", err)
			}
			if fetched.ID != id {
				t.Errorf("GetByIDAndOrg: ID = %s, want %s", fetched.ID, id)
			}
			if fetched.OrganizationID != testOrg.ID {
				t.Errorf("GetByIDAndOrg: OrganizationID = %s, want %s", fetched.OrganizationID, testOrg.ID)
			}
			if fetched.Shashu != shashu {
				t.Errorf("GetByIDAndOrg: Shashu = %s, want %s", fetched.Shashu, shashu)
			}
			fmt.Printf("✓ GetByIDAndOrg: ID=%s, Shashu=%s\n", fetched.ID, fetched.Shashu)
		})

		// 3. Update
		t.Run("Update", func(t *testing.T) {
			newID4 := "5678"
			newShashu := "Honda"
			newName := "Updated Car"
			newSekisai := 4.0

			updated, err := repo.Update(ctx, id, testOrg.ID, newID4, newShashu, &newName, &nameR, &newSekisai, &regDate, &parchDate, nil, &bumonCodeID, &driverID)
			if err != nil {
				t.Fatalf("Update failed: %v", err)
			}
			if updated.ID4 != newID4 {
				t.Errorf("Update: ID4 = %s, want %s", updated.ID4, newID4)
			}
			if updated.Shashu != newShashu {
				t.Errorf("Update: Shashu = %s, want %s", updated.Shashu, newShashu)
			}
			if updated.Name == nil || *updated.Name != newName {
				t.Errorf("Update: Name = %v, want %s", updated.Name, newName)
			}
			if updated.Sekisai == nil || *updated.Sekisai != newSekisai {
				t.Errorf("Update: Sekisai = %v, want %f", updated.Sekisai, newSekisai)
			}
			fmt.Printf("✓ Update: ID=%s, ID4=%s, Shashu=%s, Name=%s\n", updated.ID, updated.ID4, updated.Shashu, *updated.Name)
		})

		// 4. ListByOrganization
		t.Run("ListByOrganization", func(t *testing.T) {
			cars, err := repo.ListByOrganization(ctx, testOrg.ID, 10, 0)
			if err != nil {
				t.Fatalf("ListByOrganization failed: %v", err)
			}
			if len(cars) == 0 {
				t.Error("ListByOrganization: should return at least one car")
			}
			fmt.Printf("✓ ListByOrganization: returned %d cars\n", len(cars))
		})

		// 5. List
		t.Run("List", func(t *testing.T) {
			cars, err := repo.List(ctx, 10, 0)
			if err != nil {
				t.Fatalf("List failed: %v", err)
			}
			if len(cars) == 0 {
				t.Error("List: should return at least one car")
			}
			fmt.Printf("✓ List: returned %d cars\n", len(cars))
		})

		// 6. Delete
		t.Run("Delete", func(t *testing.T) {
			err := repo.Delete(ctx, id, testOrg.ID)
			if err != nil {
				t.Fatalf("Delete failed: %v", err)
			}

			// Verify hard delete
			_, err = repo.GetByIDAndOrg(ctx, id, testOrg.ID)
			if err != ErrIchibanCarNotFound {
				t.Errorf("Delete: GetByIDAndOrg should return ErrIchibanCarNotFound, got %v", err)
			}
			fmt.Printf("✓ Delete: car deleted\n")
		})
	})
}

func TestIntegration_IchibanCars_MultipleOrganizations(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewIchibanCarRepository(pool)
	orgRepo := NewOrganizationRepository(pool)
	ctx := context.Background()

	// Setup: Create two test organizations
	slug1 := fmt.Sprintf("test-ichiban-org1-%s", uuid.New().String()[:8])
	org1, err := orgRepo.Create(ctx, "Org 1", slug1)
	if err != nil {
		t.Fatalf("Setup: failed to create org1: %v", err)
	}
	defer orgRepo.Delete(ctx, org1.ID)

	slug2 := fmt.Sprintf("test-ichiban-org2-%s", uuid.New().String()[:8])
	org2, err := orgRepo.Create(ctx, "Org 2", slug2)
	if err != nil {
		t.Fatalf("Setup: failed to create org2: %v", err)
	}
	defer orgRepo.Delete(ctx, org2.ID)

	t.Run("CompositeKeyIsolation", func(t *testing.T) {
		// Same car ID can exist in different organizations
		carID := "CAR999"
		id4 := "9999"
		shashu := "Nissan"
		sekisai := 3.0

		// Create car in org1
		car1, err := repo.Create(ctx, carID, org1.ID, id4, shashu, nil, nil, &sekisai, nil, nil, nil, nil, nil)
		if err != nil {
			t.Fatalf("Create car in org1 failed: %v", err)
		}
		defer repo.Delete(ctx, carID, org1.ID)

		// Create car with same ID in org2
		car2, err := repo.Create(ctx, carID, org2.ID, id4, shashu, nil, nil, &sekisai, nil, nil, nil, nil, nil)
		if err != nil {
			t.Fatalf("Create car in org2 failed: %v", err)
		}
		defer repo.Delete(ctx, carID, org2.ID)

		// Verify both exist independently
		fetchedCar1, err := repo.GetByIDAndOrg(ctx, carID, org1.ID)
		if err != nil {
			t.Fatalf("GetByIDAndOrg for org1 failed: %v", err)
		}
		if fetchedCar1.OrganizationID != org1.ID {
			t.Errorf("Car1: OrganizationID = %s, want %s", fetchedCar1.OrganizationID, org1.ID)
		}

		fetchedCar2, err := repo.GetByIDAndOrg(ctx, carID, org2.ID)
		if err != nil {
			t.Fatalf("GetByIDAndOrg for org2 failed: %v", err)
		}
		if fetchedCar2.OrganizationID != org2.ID {
			t.Errorf("Car2: OrganizationID = %s, want %s", fetchedCar2.OrganizationID, org2.ID)
		}

		fmt.Printf("✓ CompositeKeyIsolation: Same car ID exists in both org1=%s and org2=%s\n", car1.OrganizationID, car2.OrganizationID)

		// Verify ListByOrganization returns only the org's cars
		org1Cars, err := repo.ListByOrganization(ctx, org1.ID, 100, 0)
		if err != nil {
			t.Fatalf("ListByOrganization for org1 failed: %v", err)
		}
		for _, car := range org1Cars {
			if car.OrganizationID != org1.ID {
				t.Errorf("ListByOrganization for org1 returned car from wrong org: %s", car.OrganizationID)
			}
		}
		fmt.Printf("✓ ListByOrganization: Org1 has %d cars (isolated)\n", len(org1Cars))
	})
}
