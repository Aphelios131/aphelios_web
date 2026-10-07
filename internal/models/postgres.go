package models

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type PostGresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func (cfg *PostGresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
}

func DefaultPostgresConfig() PostGresConfig {
	return PostGresConfig{
		Host:     "47.122.112.40",
		Port:     "5434",
		User:     "aphelios",
		Password: "19980131",
		Database: "aphelios_web",
		SSLMode:  "disable",
	}
}


func Open(config PostGresConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", config.String())
	if err != nil {
		return nil, fmt.Errorf("open:%w", err)
	}

	return db, nil
}