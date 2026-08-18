package customerportal

import (
	"context"
	"database/sql"
	"errors"
)

var ErrCustomerNotFound = errors.New("customer not found")

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
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
	if err := r.DB.QueryRowContext(ctx, query, userID).Scan(&customerID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrCustomerNotFound
		}
		return "", err
	}

	return customerID, nil
}
