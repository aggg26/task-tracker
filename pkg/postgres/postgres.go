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
	//db, err := pgx.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
	//	cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DatabaseName, cfg.SslMode))
	db, err := sql.Open("pgx", fmt.Sprintf("user=%s  dbname=%s password=%s host=%s port=%s sslmode=%s"))
	if err != nil {
		return nil, err
	}
	return db, nil
}
