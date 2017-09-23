package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"                         // register Postgress driver
	_ "github.com/mattes/migrate/driver/postgres" // register Migrate Postgres driver
	"github.com/mattes/migrate/migrate"
	"github.com/onlyafly/tradgard/pkg/server"
)

const (
	defaultPort = "5000"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, "[FATAL] Panic recovered in main", r)
		}
	}()

	// DATABASE_URL should look like: postgres://127.0.0.1:5432/tradgard?sslmode=disable
	dbURLStr := os.Getenv("DATABASE_URL")
	dbURL, err := url.Parse(dbURLStr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[FATAL] Failed to parse envvar DATABASE_URL", dbURLStr, err)
		return
	}

	fmt.Println("[INFO] Connecting to DB", dbURL.Host)
	db, err := connect(*dbURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[FATAL] Failed to connect to DB", err)
		return
	}
	defer db.Close()

	fmt.Println("[INFO] Migrating DB", dbURL.Host)
	if err = migrateDB(dbURL, "./etc/db/"); err != nil {
		fmt.Fprintln(os.Stderr, "[FATAL] Failed to migrate DB", err)
		return
	}

	config := server.Config{
		Port:     getEnvOr("PORT", defaultPort),
		SiteName: getEnvOr("SITE_NAME", "My Garden"),
		Database: db,
	}

	server.Start(config)
}

// connect establishes connection to a database and probes connectivity.
func connect(u url.URL) (*sqlx.DB, error) {
	dsn := u.String()

	db, err := sqlx.Connect(u.Scheme, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN or no database found: %v", err)
	}

	return db, nil
}

func migrateDB(dbURL *url.URL, migrationsPath string) error {
	migrate.Graceful() // Set graceful handling of ^C by migrate
	if migrateErrs, ok := migrate.UpSync(dbURL.String(), migrationsPath); !ok {
		return fmt.Errorf("%v", migrateErrs)
	}
	return nil
}

func getEnvOr(envVar, fallback string) string {
	if result := os.Getenv(envVar); result != "" {
		return result
	}
	return fallback
}
