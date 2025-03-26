package config

import (
	"cmp"
	"flag"
	"os"
	"strconv"
)

type Config struct {
	Host        string
	Port        int
	DbDsn       string
	MigratePath string
	Debug       bool
}

const (
	defaultPort        = 8080
	defaultHost        = "localhost"
	defaultDbDst       = "postgres://user:password@localhost:5432/gt5?sslmode=disable"
	defaultMigratePath = "migrations"
)

func ReadConfig() (*Config, error) {
	var cfg Config

	flag.StringVar(&cfg.Host, "host", defaultHost, "flag for explicit server host specifications")
	flag.IntVar(&cfg.Port, "port", defaultPort, "flag for explicit server port specifications")
	flag.StringVar(&cfg.DbDsn, "db", defaultDbDst, "flag for explicit db connection string")
	flag.StringVar(&cfg.MigratePath, "migrate", defaultMigratePath, "flag for explicit migrate path")
	flag.BoolVar(&cfg.Debug, "debug", false, "flag for explicit debug mode")

	flag.Parse()

	if cfg.Host == "localhost" {
		cfg.Host = cmp.Or(os.Getenv("HOST"), cfg.Host)
	}
	if cfg.Port == 8080 {
		defPort := strconv.Itoa(cfg.Port)
		envPort := cmp.Or(os.Getenv("PORT"), defPort)
		port, err := strconv.Atoi(envPort)
		if err != nil {
			return nil, err
		}
		cfg.Port = port
	}
	if cfg.DbDsn == defaultDbDst {
		cfg.DbDsn = cmp.Or(os.Getenv("DB_DSN"), cfg.DbDsn)
	}
	if cfg.MigratePath == defaultMigratePath {
		cfg.MigratePath = cmp.Or(os.Getenv("MIGRATE_PATH"), cfg.MigratePath)
	}

	return &cfg, nil
}
