package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"jedug-backend/internal/config"
)

var Pool *pgxpool.Pool

func Connect(cfg config.DBConfig) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	poolCfg, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	poolCfg.MaxConns = 10
	poolCfg.MinConns = 2

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}

	Pool = pool
	log.Println("✅ Connected to PostgreSQL")
	return pool, nil
}

func TryConnect(cfg config.DBConfig) *pgxpool.Pool {
	pool, err := Connect(cfg)
	if err != nil {
		log.Printf("⚠️  Database not available: %v (server will start without DB)", err)
		return nil
	}
	return pool
}

func Close() {
	if Pool != nil {
		Pool.Close()
		log.Println("🔌 PostgreSQL connection closed")
	}
}
