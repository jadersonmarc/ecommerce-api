package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	defaultHost          = "localhost"
	defaultPort          = "5432"
	defaultUser          = "postgres"
	defaultPassword      = "123mudar"
	defaultDatabase      = "ecommerce"
	defaultSSLMode       = "disable"
	defaultConnectTimout = 5 * time.Second
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func LoadConfig() Config {
	return Config{
		Host:     envOrDefault("POSTGRES_HOST", defaultHost),
		Port:     envOrDefault("POSTGRES_PORT", defaultPort),
		User:     envOrDefault("POSTGRES_USER", defaultUser),
		Password: envOrDefault("POSTGRES_PASSWORD", defaultPassword),
		Name:     envOrDefault("POSTGRES_DB", defaultDatabase),
		SSLMode:  envOrDefault("POSTGRES_SSLMODE", defaultSSLMode),
	}
}

func (c Config) dsn(databaseName string) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		databaseName,
		c.SSLMode,
	)
}

func Open(ctx context.Context, cfg Config) (*gorm.DB, error) {
	if err := ensureDatabase(ctx, cfg); err != nil {
		return nil, err
	}

	gormDB, err := gorm.Open(postgres.Open(cfg.dsn(cfg.Name)), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db, err := gormDB.DB()
	if err != nil {
		return nil, err
	}

	pingCtx, cancel := context.WithTimeout(ctx, defaultConnectTimout)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		db.Close()
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err := migrate(ctx, gormDB); err != nil {
		db.Close()
		return nil, err
	}

	return gormDB, nil
}

func ensureDatabase(ctx context.Context, cfg Config) error {
	adminDB, err := sql.Open("pgx", cfg.dsn("postgres"))
	if err != nil {
		return err
	}
	defer adminDB.Close()

	pingCtx, cancel := context.WithTimeout(ctx, defaultConnectTimout)
	defer cancel()

	if err := adminDB.PingContext(pingCtx); err != nil {
		return err
	}

	var exists bool
	if err := adminDB.QueryRowContext(
		ctx,
		`SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)`,
		cfg.Name,
	).Scan(&exists); err != nil {
		return err
	}

	if exists {
		return nil
	}

	_, err = adminDB.ExecContext(ctx, fmt.Sprintf(`CREATE DATABASE "%s"`, cfg.Name))
	if err != nil {
		return err
	}

	return nil
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

var ErrNotFound = errors.New("not found")
