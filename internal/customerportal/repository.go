package customerportal

import (
	"context"
	"database/sql"
	"errors"
)

var ErrCustomerNotFound = errors.New("customer not found")

type DB interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

type Repository struct {
	DB DB
}

func NewRepository(db DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) CustomerIDForUser(
	ctx context.Context,
	userID string,
) (string, error) {
	if r == nil || r.DB == nil {
		return "", errors.New("customer portal repository is not configured")
	}

	const query = `
SELECT customer_id
FROM apb.customer_identities
WHERE user_id = $1
`

	var customerID string

	if err := r.DB.QueryRowContext(
		ctx,
		query,
		userID,
	).Scan(&customerID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrCustomerNotFound
		}

		return "", err
	}

	return customerID, nil
}

func (r *Repository) GetCustomer(
	ctx context.Context,
	customerID string,
) (CustomerProfile, error) {
	if r == nil || r.DB == nil {
		return CustomerProfile{}, errors.New("customer portal repository is not configured")
	}

	const query = `
SELECT
	id,
	customer_code,
	name,
	COALESCE(email, ''),
	COALESCE(phone, ''),
	COALESCE(address, ''),
	COALESCE(village, ''),
	COALESCE(district, ''),
	COALESCE(regency, ''),
	COALESCE(province, ''),
	COALESCE(postal_code, ''),
	status
FROM apb.customers
WHERE id = $1
`

	var c CustomerProfile

	err := r.DB.QueryRowContext(
		ctx,
		query,
		customerID,
	).Scan(
		&c.ID,
		&c.Code,
		&c.Name,
		&c.Email,
		&c.Phone,
		&c.Address,
		&c.Village,
		&c.District,
		&c.Regency,
		&c.Province,
		&c.PostalCode,
		&c.Status,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return CustomerProfile{}, ErrCustomerNotFound
	}

	if err != nil {
		return CustomerProfile{}, err
	}

	return c, nil
}

func (r *Repository) GetServices(
	ctx context.Context,
	customerID string,
) ([]CustomerService, error) {
	if r == nil || r.DB == nil {
		return nil, errors.New("customer portal repository is not configured")
	}

	const query = `
SELECT
	id,
	service_code,
	service_type,
	COALESCE(package_name, ''),
	COALESCE(download_speed, 0),
	COALESCE(upload_speed, 0),
	status
FROM apb.services
WHERE customer_id = $1
ORDER BY created_at DESC
`

	rows, err := r.DB.QueryContext(
		ctx,
		query,
		customerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	services := make([]CustomerService, 0)

	for rows.Next() {
		var s CustomerService

		if err := rows.Scan(
			&s.ID,
			&s.ServiceCode,
			&s.ServiceType,
			&s.PackageName,
			&s.DownloadSpeed,
			&s.UploadSpeed,
			&s.Status,
		); err != nil {
			return nil, err
		}

		services = append(services, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return services, nil
}
