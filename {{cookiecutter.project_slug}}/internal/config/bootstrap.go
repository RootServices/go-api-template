package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"{{cookiecutter.module_name}}/internal/gcp"
)

// structs to help fetching secrets from gcp

type SecretCoordinates struct {
	ProjectNumber string
	DBPasswordKey string
}

type Secrets struct {
	DBPassword string
}

// structs returned by load fn
type Database struct {
	DSN string // Data Source Name Native Postgres
}

type AppConfig struct {
	Env                   string
	DB                    Database
	ProjectID             string
	StorageBucket         string
	StorageServiceAccount string
}

type GetVariable func(key string) string

func readVariable(key string) string {
	return os.Getenv(key)
}

type bootStrap struct {
	getVariable GetVariable
	repo        gcp.SecretRepository
	log         *slog.Logger
}

type BootStrap interface {
	Load(ctx context.Context) (*AppConfig, error)
	FetchSecrets(ctx context.Context, coords SecretCoordinates) (Secrets, error)
}

func BootStrapFactory(ctx context.Context, log *slog.Logger) (BootStrap, error) {
	getVariable := readVariable

	if getVariable("ENV") == "local" {
		return NewLocalBootStrap(context.Background(), log)
	}

	return NewBootStrap(context.Background(), log)
}

func NewBootStrap(ctx context.Context, log *slog.Logger) (BootStrap, error) {
	repo, err := gcp.NewSecretRepository(ctx, log)
	if err != nil {
		return nil, fmt.Errorf("failed to create secret repository: %w", err)
	}

	return &bootStrap{
		getVariable: readVariable,
		repo:        repo,
		log:         log,
	}, nil
}

func NewLocalBootStrap(ctx context.Context, log *slog.Logger) (BootStrap, error) {
	return &bootStrap{
		getVariable: readVariable,
		repo:        gcp.NewFakeSecretRepo(),
		log:         log,
	}, nil
}

func (b *bootStrap) Load(ctx context.Context) (*AppConfig, error) {

	env := b.getVariable("ENV")

	b.log.Info(fmt.Sprintf("using config %s", env))
	var dsn string
	dsnTemplate := "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s"

	dbHost := b.getVariable("DB_HOST")
	dbName := b.getVariable("DB_NAME")
	dbPort := b.getVariable("DB_PORT")
	dbSSLMode := b.getVariable("DB_SSL_MODE")
	dbUser := b.getVariable("DB_USER")

	if env == "local" {
		dbPassword := b.getVariable("DB_PASSWORD")
		dsn = fmt.Sprintf(dsnTemplate, dbHost, dbUser, dbPassword, dbName, dbPort, dbSSLMode)
	} else {
		// get secrts from gcp
		gcpProjectNumber := b.getVariable("GCP_PROJECT_NUMBER")
		dbPasswordKey := b.getVariable("DB_PASSWORD_KEY")

		coords := SecretCoordinates{
			ProjectNumber: gcpProjectNumber,
			DBPasswordKey: dbPasswordKey,
		}

		secrets, err := b.FetchSecrets(ctx, coords)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch secrets: %w", err)
		}

		dsn = fmt.Sprintf(dsnTemplate, dbHost, dbUser, secrets.DBPassword, dbName, dbPort, dbSSLMode)
	}

	storageBucket := b.getVariable("STORAGE_BUCKET")

	// 3. Populate AppConfig
	appConfig := &AppConfig{
		Env: env,
		DB: Database{
			DSN: dsn,
		},
		ProjectID:             b.getVariable("GCP_PROJECT_ID"),
		StorageBucket:         storageBucket,
		StorageServiceAccount: b.getVariable("STORAGE_SERVICE_ACCOUNT"),
	}

	return appConfig, nil
}

// fetches secrets from gcp
func (b *bootStrap) FetchSecrets(ctx context.Context, coords SecretCoordinates) (Secrets, error) {
	secrets := Secrets{}

	if coords.DBPasswordKey != "" && coords.ProjectNumber != "" {
		val, err := b.repo.GetSecret(ctx, coords.ProjectNumber, coords.DBPasswordKey, "latest")
		if err != nil {
			return Secrets{}, fmt.Errorf("failed to fetch secret 'dbPassword' (project: %s, secret: %s, version: latest): %w", coords.ProjectNumber, coords.DBPasswordKey, err)
		}
		secrets.DBPassword = val
	}

	return secrets, nil
}
