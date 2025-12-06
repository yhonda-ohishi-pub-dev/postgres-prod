package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrIchibanCarNotFound = errors.New("ichiban car not found")
)

// IchibanCar represents the database model
type IchibanCar struct {
	ID             string
	OrganizationID string
	ID4            string
	Name           *string
	NameR          *string
	Shashu         string
	Sekisai        *float64
	RegDate        *string
	ParchDate      *string
	ScrapDate      *string
	BumonCodeID    *string
	DriverID       *string
}

// IchibanCarRepository handles database operations for ichiban_cars
type IchibanCarRepository struct {
	db DB
}

// NewIchibanCarRepository creates a new repository
func NewIchibanCarRepository(pool *pgxpool.Pool) *IchibanCarRepository {
	return &IchibanCarRepository{db: pool}
}

// NewIchibanCarRepositoryWithDB creates a repository with custom DB interface (for testing)
func NewIchibanCarRepositoryWithDB(db DB) *IchibanCarRepository {
	return &IchibanCarRepository{db: db}
}

// Create inserts a new ichiban car
func (r *IchibanCarRepository) Create(ctx context.Context, id, organizationID, id4, shashu string, name, nameR *string, sekisai *float64, regDate, parchDate, scrapDate, bumonCodeID, driverID *string) (*IchibanCar, error) {
	query := `
		INSERT INTO ichiban_cars (id, organization_id, id4, name, "name_R", shashu, sekisai, reg_date, parch_date, scrap_date, bumon_code_id, driver_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, organization_id, id4, name, "name_R", shashu, sekisai, reg_date, parch_date, scrap_date, bumon_code_id, driver_id
	`

	var car IchibanCar
	err := r.db.QueryRow(ctx, query, id, organizationID, id4, name, nameR, shashu, sekisai, regDate, parchDate, scrapDate, bumonCodeID, driverID).Scan(
		&car.ID, &car.OrganizationID, &car.ID4, &car.Name, &car.NameR, &car.Shashu, &car.Sekisai, &car.RegDate, &car.ParchDate, &car.ScrapDate, &car.BumonCodeID, &car.DriverID,
	)
	if err != nil {
		return nil, err
	}

	return &car, nil
}

// GetByIDAndOrg retrieves an ichiban car by composite primary key (id, organization_id)
func (r *IchibanCarRepository) GetByIDAndOrg(ctx context.Context, id, organizationID string) (*IchibanCar, error) {
	query := `
		SELECT id, organization_id, id4, name, "name_R", shashu, sekisai, reg_date, parch_date, scrap_date, bumon_code_id, driver_id
		FROM ichiban_cars
		WHERE id = $1 AND organization_id = $2
	`

	var car IchibanCar
	err := r.db.QueryRow(ctx, query, id, organizationID).Scan(
		&car.ID, &car.OrganizationID, &car.ID4, &car.Name, &car.NameR, &car.Shashu, &car.Sekisai, &car.RegDate, &car.ParchDate, &car.ScrapDate, &car.BumonCodeID, &car.DriverID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrIchibanCarNotFound
		}
		return nil, err
	}

	return &car, nil
}

// Update modifies an existing ichiban car
func (r *IchibanCarRepository) Update(ctx context.Context, id, organizationID, id4, shashu string, name, nameR *string, sekisai *float64, regDate, parchDate, scrapDate, bumonCodeID, driverID *string) (*IchibanCar, error) {
	query := `
		UPDATE ichiban_cars
		SET id4 = $3, name = $4, "name_R" = $5, shashu = $6, sekisai = $7, reg_date = $8, parch_date = $9, scrap_date = $10, bumon_code_id = $11, driver_id = $12
		WHERE id = $1 AND organization_id = $2
		RETURNING id, organization_id, id4, name, "name_R", shashu, sekisai, reg_date, parch_date, scrap_date, bumon_code_id, driver_id
	`

	var car IchibanCar
	err := r.db.QueryRow(ctx, query, id, organizationID, id4, name, nameR, shashu, sekisai, regDate, parchDate, scrapDate, bumonCodeID, driverID).Scan(
		&car.ID, &car.OrganizationID, &car.ID4, &car.Name, &car.NameR, &car.Shashu, &car.Sekisai, &car.RegDate, &car.ParchDate, &car.ScrapDate, &car.BumonCodeID, &car.DriverID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrIchibanCarNotFound
		}
		return nil, err
	}

	return &car, nil
}

// Delete hard-deletes an ichiban car (no soft delete for this table)
func (r *IchibanCarRepository) Delete(ctx context.Context, id, organizationID string) error {
	query := `
		DELETE FROM ichiban_cars
		WHERE id = $1 AND organization_id = $2
	`

	result, err := r.db.Exec(ctx, query, id, organizationID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrIchibanCarNotFound
	}

	return nil
}

// ListByOrganization retrieves ichiban cars by organization with pagination
func (r *IchibanCarRepository) ListByOrganization(ctx context.Context, organizationID string, limit int, offset int) ([]*IchibanCar, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT id, organization_id, id4, name, "name_R", shashu, sekisai, reg_date, parch_date, scrap_date, bumon_code_id, driver_id
		FROM ichiban_cars
		WHERE organization_id = $1
		ORDER BY id
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cars []*IchibanCar
	for rows.Next() {
		var car IchibanCar
		err := rows.Scan(&car.ID, &car.OrganizationID, &car.ID4, &car.Name, &car.NameR, &car.Shashu, &car.Sekisai, &car.RegDate, &car.ParchDate, &car.ScrapDate, &car.BumonCodeID, &car.DriverID)
		if err != nil {
			return nil, err
		}
		cars = append(cars, &car)
	}

	return cars, rows.Err()
}

// List retrieves all ichiban cars with pagination
func (r *IchibanCarRepository) List(ctx context.Context, limit int, offset int) ([]*IchibanCar, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query := `
		SELECT id, organization_id, id4, name, "name_R", shashu, sekisai, reg_date, parch_date, scrap_date, bumon_code_id, driver_id
		FROM ichiban_cars
		ORDER BY organization_id, id
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cars []*IchibanCar
	for rows.Next() {
		var car IchibanCar
		err := rows.Scan(&car.ID, &car.OrganizationID, &car.ID4, &car.Name, &car.NameR, &car.Shashu, &car.Sekisai, &car.RegDate, &car.ParchDate, &car.ScrapDate, &car.BumonCodeID, &car.DriverID)
		if err != nil {
			return nil, err
		}
		cars = append(cars, &car)
	}

	return cars, rows.Err()
}
