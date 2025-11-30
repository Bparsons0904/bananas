package repositories

import (
	"bananas/internal/database"
	"bananas/internal/logger"
)

type Manager struct {
	repos  map[string]RepositoryInterface
	Logger logger.Logger
}

func NewManager(db *database.DB) (*Manager, error) {
	log := logger.New("repository-manager")
	repos := make(map[string]RepositoryInterface)

	repos["sql"] = NewSQLRepository(db)
	repos["gorm"] = NewGORMRepository(db)
	repos["sqlx"] = NewSQLxRepository(db)
	repos["pgx"] = NewPGXRepository(db)

	log.Info("Initialized all ORM repositories: sql, gorm, sqlx, pgx")

	return &Manager{
		repos:  repos,
		Logger: log,
	}, nil
}

func (m *Manager) GetRepository(ormType string) RepositoryInterface {
	if ormType == "" {
		ormType = "sql"
	}

	if repo, ok := m.repos[ormType]; ok {
		return repo
	}

	m.Logger.Info("Unknown ORM type '%s', falling back to 'sql'", ormType)
	return m.repos["sql"]
}

func (m *Manager) ListAvailableORMs() []string {
	return []string{"sql", "gorm", "sqlx", "pgx"}
}
