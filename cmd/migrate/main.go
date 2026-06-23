package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/tjahdja/SimpleBlog/internal/database/migrations"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/golang-migrate/migrate/v4/source/iofs"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	sourceDriver, err := iofs.New(migrations.Files, ".")
	if err != nil {
		log.Fatalf("Failed to initialize migration files source: %v", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", sourceDriver, dsn)
	if err != nil {
		log.Fatalf("Failed to connect to migration database: %v", err)
	}

	action := flag.String("action", "up", "Migration action to execute: 'up' or 'down'")
	flag.Parse()

	switch *action {
	case "up":
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("Migration UP error: %v", err)
		}
		fmt.Println("Database schemas migrated UP successfully!")
	case "down":
		if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("Migration DOWN error: %v", err)
		}
		fmt.Println("Database schemas rolled DOWN cleanly!")
	default:
		log.Fatalf("Unknown migration flag option choice: %s. Use 'up' or 'down'.", *action)
	}
}
