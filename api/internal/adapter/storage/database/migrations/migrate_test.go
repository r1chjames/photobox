package migrations

import (
	"strings"
	"testing"
)

func TestToMigrateURL(t *testing.T) {
	dsn := "host=db.example.com user=photobox password=p@ss word dbname=photoboxdb port=5433 sslmode=disable search_path=photobox,public"
	got, err := toMigrateURL(dsn)
	if err != nil {
		t.Fatalf("toMigrateURL returned error: %v", err)
	}

	wantScheme := "postgres://"
	if !strings.HasPrefix(got, wantScheme) {
		t.Errorf("expected scheme %q, got %q", wantScheme, got)
	}
	for _, want := range []string{
		"db.example.com:5433",
		"/photoboxdb",
		"sslmode=disable",
		"search_path=photobox%2Cpublic",
		"x-multi-statement=true",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("migrate URL %q missing %q", got, want)
		}
	}
	// Password is URL-escaped so it survives as a single query/user component.
	if strings.Contains(got, "p@ss word") {
		t.Errorf("password not escaped in migrate URL: %q", got)
	}
}

func TestToMigrateURL_RequiresDbName(t *testing.T) {
	if _, err := toMigrateURL("host=localhost user=photobox sslmode=disable"); err == nil {
		t.Error("expected error when dbname missing from DSN")
	}
}

func TestToMigrateURL_Defaults(t *testing.T) {
	got, err := toMigrateURL("host=h user=u password=p dbname=d")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{"h:5432", "/d", "postgres://"} {
		if !strings.Contains(got, want) {
			t.Errorf("migrate URL %q missing default %q", got, want)
		}
	}
}
