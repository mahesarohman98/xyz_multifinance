package testutils

import (
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
)

func SeedCustomer(t *testing.T, db *sqlx.DB, id, name string, nik string) {
	t.Helper()
	t.Helper()

	const query = `
	INSERT INTO Customers (
		customer_id,
		nik,
		full_name,
		legal_name,
		place_of_birth,
		date_of_birth,
		wages,
		ktp_photo_url,
		photo_url,
		created_at,
		updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := db.Exec(query,
		id,
		nik,
		name,
		name,
		"tangerang",
		"1998-04-01",
		3500000,
		"https://imgurl.com/ktp",
		"https://imgurl.com/photo",
		time.Date(2025, 4, 1, 1, 0, 0, 0, time.UTC),
		time.Date(2025, 4, 1, 1, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("failed to seed customer: %v", err)
	}
}

func SeedSources(t *testing.T, db *sqlx.DB) {
	t.Helper()
	t.Helper()

	const query = `
	INSERT INTO Sources (
		source_id,
		category,
		name
	) VALUES (?, ?, ?), (?, ?, ?)
	`

	_, err := db.Exec(query,
		"sourceid-1",
		"ecommerce",
		"tokopakedi",

		"sourceid-2",
		"dealer",
		"dealer name",
	)
	if err != nil {
		t.Fatalf("failed to seed source: %v", err)
	}
}
