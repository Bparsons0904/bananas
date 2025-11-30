package database

import (
	"bananas/internal/config"
	"bananas/internal/logger"
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type DB struct {
	SQL    *sql.DB
	GORM   *gorm.DB
	SQLx   *sqlx.DB
	PGX    *pgxpool.Pool
	Logger logger.Logger
	Config config.Config
}

func New(cfg config.Config) (*DB, error) {
	log := logger.New("database")
	dsn := cfg.GetDatabaseDSN()

	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Er("failed to open database/sql connection", err)
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		log.Er("failed to ping database/sql", err)
		return nil, err
	}
	log.Info("database/sql connection established")

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		log.Er("failed to open GORM connection", err)
		return nil, err
	}
	log.Info("GORM connection established")

	sqlxDB, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Er("failed to open SQLx connection", err)
		return nil, err
	}
	log.Info("SQLx connection established")

	pgxPool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Er("failed to create PGX pool", err)
		return nil, err
	}
	if err := pgxPool.Ping(context.Background()); err != nil {
		log.Er("failed to ping PGX pool", err)
		return nil, err
	}
	log.Info("PGX pool established")

	log.Info("All database connections established")

	return &DB{
		SQL:    sqlDB,
		GORM:   gormDB,
		SQLx:   sqlxDB,
		PGX:    pgxPool,
		Logger: log,
		Config: cfg,
	}, nil
}

func (db *DB) Close() error {
	var lastErr error

	if db.SQL != nil {
		if err := db.SQL.Close(); err != nil {
			db.Logger.Er("failed to close database/sql", err)
			lastErr = err
		}
	}

	if db.SQLx != nil {
		if err := db.SQLx.Close(); err != nil {
			db.Logger.Er("failed to close SQLx", err)
			lastErr = err
		}
	}

	if db.PGX != nil {
		db.PGX.Close()
	}

	return lastErr
}

func (db *DB) HealthCheck(ctx context.Context) error {
	return db.SQL.PingContext(ctx)
}