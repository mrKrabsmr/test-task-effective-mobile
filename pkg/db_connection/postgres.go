package dbConnection

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mrKrabsmr/commerce-edu-api/internal/configs"
)

type PostgresDB struct {
	config *configs.Config
}

func NewPGConnection(config *configs.Config) *PostgresDB {
	return &PostgresDB{
		config: config,
	}
}

func (p *PostgresDB) PostgreSQLConnection() (*sqlx.DB, error) {
	db, err := sqlx.Connect(p.config.DBDialect, p.config.DBAddress)
	if err != nil {
		return nil, fmt.Errorf("error, not connected to database, %w", err)
	}

	if err = db.Ping(); err != nil {
		defer db.Close()
		return nil, fmt.Errorf("error, not sent ping to database, %w", err)
	}

	return db, nil
}
