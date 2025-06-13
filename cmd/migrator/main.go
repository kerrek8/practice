package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	var migrationsTable, migrationsPath string

	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")
	flag.Parse()
	if migrationsPath == "" && migrationsTable == "" {
		fmt.Println("Please provide migrations-path flags")
		return
	}

	m, err := migrate.New("file://"+migrationsPath, fmt.Sprintf("sqlite:/../database/database.db?x-migrations-table=%s", migrationsTable))
	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no changes")
			return
		}
		panic(err)
	}
}
