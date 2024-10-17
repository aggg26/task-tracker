package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresConfig struct {
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	DatabaseName string `yaml:"dbname"`
	SslMode      string `yaml:"sslMode"`
}

func NewPostgresDb(cfg PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=%s",
		cfg.Username, cfg.DatabaseName, cfg.Password, cfg.Host, cfg.Port, cfg.SslMode))
	if err != nil {
		return nil, err
	}
	return db, nil
}
