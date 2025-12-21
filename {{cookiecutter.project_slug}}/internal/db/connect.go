package db

import (
	"log/slog"
	"os"

	"cloud.google.com/go/cloudsqlconn/postgres/pgxv5"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type MakeDbFn func(dsn string, log *slog.Logger) *gorm.DB

func MakeDbFactory(env string) MakeDbFn {
	if env == "local" {
		return MakeLocalDb
	}
	return MakeCloudSQLDb
}

func MakeLocalDb(dsn string, log *slog.Logger) *gorm.DB {
	log.Info("connecting to local postgresdb")
	// Initialize database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
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

// this uses the cloud sql connector approach
// https://github.com/go-gorm/gorm/issues/6991
func MakeCloudSQLDb(dsn string, log *slog.Logger) *gorm.DB {
	log.Info("connecting to cloud sql")
	cleanup, err := pgxv5.RegisterDriver("cloudsql-postgres")
	if err != nil {
		panic(err)
	}
	defer cleanup()

	db, err := gorm.Open(postgres.New(postgres.Config{
		DriverName: "cloudsql-postgres",
		DSN:        dsn,
	}), &gorm.Config{
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
