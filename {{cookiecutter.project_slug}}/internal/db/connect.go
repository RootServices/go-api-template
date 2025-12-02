package db

import (
	"log/slog"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func MakeDb(databaseURL string, log *slog.Logger) *gorm.DB {
	// Initialize database connection
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Error("failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	return db
}

func MakeDbSqlite() (*gorm.DB, error) {
	// Initialize database connection
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	return db, err
}
