package postgresql

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Govorov1705/avito-pvz/configs"
	"github.com/Govorov1705/avito-pvz/pkg/logger"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	Pool *pgxpool.Pool
}

func NewStorage() *Storage {
	dbpool, err := pgxpool.New(context.Background(), configs.Cfg.DBURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	logger.Logger.Info("DB pool created")

	// Get absolute path to migrations
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		logger.Logger.Fatal("Could not determine caller information for migrations")
	}
	migrationsPath := filepath.Join(filepath.Dir(filename), "migrations")
	absMigrationsPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		logger.Logger.Fatal("Error converting migrations path to absolute", zap.Error(err))
	}
	migrationSourceURL := fmt.Sprintf("file://%s", absMigrationsPath)

	m, err := migrate.New(
		migrationSourceURL,
		configs.Cfg.DBURL+"sslmode=disable")
	if err != nil {
		logger.Logger.Fatal("Error creating Migrate instance", zap.Error(err))
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Logger.Fatal("Error applying migrations", zap.Error(err))
	}
	logger.Logger.Info("Migrations applied")

	return &Storage{
		Pool: dbpool,
	}
}
