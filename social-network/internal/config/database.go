package config

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Database конфигурация базы данных.
type Database struct {
	Host     string `envconfig:"HOST" required:"true"`
	Port     int    `envconfig:"PORT" required:"true"`
	User     string `envconfig:"USER" required:"true"`
	Name     string `envconfig:"NAME" required:"true"`
	Password string `envconfig:"PASSWORD" required:"true"`
	SSLMode  string `envconfig:"SSL_MODE" default:"disable"`
	Schema   string `envconfig:"SCHEMA" default:"public"`

	ConnMaxLifeTime      time.Duration `envconfig:"CONN_MAX_LIFETIME"`
	MaxOpenConns         int32         `envconfig:"MAX_OPEN_CONNS"`
	MaxIdleConns         int32         `envconfig:"MAX_IDLE_CONNS"`
	MaxIdleConnsLifeTime time.Duration `envconfig:"MAX_IDLE_CONNS_LIFETIME" default:"30s"`
}

// URL строка подключения к бд.
func (d Database) url() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=%s search_path=%s",
		d.Host,
		d.Port,
		d.User,
		d.Name,
		d.Password,
		d.SSLMode,
		d.Schema,
	)
}

func (d Database) PgxConfig() (*pgxpool.Config, error) {
	cfg, err := pgxpool.ParseConfig(d.url())
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgreSQL's config: %w", err)
	}

	cfg.MaxConnIdleTime = d.MaxIdleConnsLifeTime
	cfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeExec

	if d.ConnMaxLifeTime > 0 {
		cfg.MaxConnLifetime = d.ConnMaxLifeTime
	}

	if d.MaxOpenConns > 0 {
		cfg.MaxConns = d.MaxOpenConns
	}

	cfg.MaxConnIdleTime = d.MaxIdleConnsLifeTime

	return cfg, nil
}

type ReplicaDatabase struct {
	Enabled bool `envconfig:"ENABLED"`
	Database
}

type MongoDatabase struct {
	Host         string `envconfig:"HOST" required:"true"`
	Port         int    `envconfig:"PORT" required:"true"`
	User         string `envconfig:"USER"`
	Password     string `envconfig:"PASSWORD"`
	Name         string `envconfig:"NAME" required:"true"`
	AuthRequired bool   `envconfig:"AUTH_REQUIRED" default:"true"`
}

func (d MongoDatabase) URI() string {
	hostPort := net.JoinHostPort(d.Host, strconv.Itoa(d.Port))

	if d.AuthRequired {
		return fmt.Sprintf("mongodb://%s:%s@%s/%s", d.User, d.Password, hostPort, d.Name)
	}

	return fmt.Sprintf("mongodb://%s/%s", hostPort, d.Name)
}

func (d MongoDatabase) Validate() error {
	if d.AuthRequired && (d.User == "" || d.Password == "") {
		return fmt.Errorf("invalid user/password")
	}

	return nil
}

type TarantoolDB struct {
	Host string `envconfig:"HOST" required:"true"`
	Port string `envconfig:"PORT" required:"true"`
	User string `envconfig:"USER" required:"true"`
	Pass string `envconfig:"PASS"`
}
