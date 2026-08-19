package migrations

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMigrationFilesAreSorted(t *testing.T) {
	dir := t.TempDir()

	files := map[string]string{
		"0003_customer.sql": "SELECT 3;",
		"0001_init.sql":     "SELECT 1;",
		"0002_auth.sql":     "SELECT 2;",
		"README.txt":        "ignored",
	}

	for name, content := range files {
		if err := os.WriteFile(
			filepath.Join(dir, name),
			[]byte(content),
			0600,
		); err != nil {
			t.Fatal(err)
		}
	}

	got, err := migrationFiles(dir)
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 3 {
		t.Fatalf("expected 3 migrations, got %d", len(got))
	}

	expected := []string{
		"0001_init",
		"0002_auth",
		"0003_customer",
	}

	for i, want := range expected {
		if got[i].Version != want {
			t.Fatalf(
				"migration %d: expected %q, got %q",
				i,
				want,
				got[i].Version,
			)
		}
	}
}

func TestMigrationFilesMissingDirectory(t *testing.T) {
	_, err := migrationFiles(
		filepath.Join(
			t.TempDir(),
			"does-not-exist",
		),
	)

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestMigrationFilesIgnoresDirectories(t *testing.T) {
	dir := t.TempDir()

	if err := os.Mkdir(
		filepath.Join(dir, "0002_directory.sql"),
		0700,
	); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(
		filepath.Join(dir, "0001_init.sql"),
		[]byte("SELECT 1;"),
		0600,
	); err != nil {
		t.Fatal(err)
	}

	got, err := migrationFiles(dir)
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 1 {
		t.Fatalf(
			"expected 1 migration, got %d",
			len(got),
		)
	}

	if got[0].Version != "0001_init" {
		t.Fatalf(
			"unexpected migration %q",
			got[0].Version,
		)
	}
}
