package customerportal

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestCustomerOwnershipIntegration(t *testing.T) {
	databaseURL := os.Getenv("A_RADIUS_TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("A_RADIUS_TEST_DATABASE_URL is not set")
	}

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()

	if err := db.PingContext(ctx); err != nil {
		t.Fatal(err)
	}

	repo := NewRepository(db)

	const userA = "fae05457-5ffa-40c4-8d61-047459c77870"

	customerA, err := repo.CustomerIDForUser(ctx, userA)
	if err != nil {
		t.Fatal(err)
	}

	if customerA != "db0b78cb-99c6-4cd2-adba-018845e62e54" {
		t.Fatalf(
			"unexpected customer mapping: got %s",
			customerA,
		)
	}

	customer, err := repo.GetCustomer(ctx, customerA)
	if err != nil {
		t.Fatal(err)
	}

	if customer.ID != customerA {
		t.Fatalf(
			"customer ownership mismatch: got %s want %s",
			customer.ID,
			customerA,
		)
	}

	services, err := repo.GetServices(ctx, customerA)
	if err != nil {
		t.Fatal(err)
	}

	for _, service := range services {
		var actualCustomerID string

		err := db.QueryRowContext(
			ctx,
			`SELECT customer_id FROM apb.services WHERE id = $1`,
			service.ID,
		).Scan(&actualCustomerID)
		if err != nil {
			t.Fatal(err)
		}

		if actualCustomerID != customerA {
			t.Fatalf(
				"service ownership violation: service=%s customer=%s expected=%s",
				service.ID,
				actualCustomerID,
				customerA,
			)
		}
	}
}
