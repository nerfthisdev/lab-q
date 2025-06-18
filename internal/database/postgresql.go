package database

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	ConnectionURI   string        `yaml:"connectionURI"`
	MaxConn         int32         `yaml:"maxConn"`
	MaxConnLifetime time.Duration `yaml:"maxConnLifetime"`
}

func NewPostgresDB(cfg Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), "postgres://nerfbot:66sTE079@94.228.124.206:5432/botdb")
	if err != nil {
		return &pgxpool.Pool{}, errors.New("failed to parse connectionURI to build pool")
	}
	if err := pool.Ping(context.Background()); err != nil {
		return &pgxpool.Pool{}, errors.New("failed to ping database")
	}
	return pool, nil
}
