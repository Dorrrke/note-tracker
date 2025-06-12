package config

import (
	"cmp"
	"flag"
	"os"
	"strconv"
)

type Config struct {
	Host  string
	Port  int
	Debug bool
	DBDSN string
	// MigratePath string
}

const (
	defaultPort  = 8080
	defaultHost  = "localhost"
	defaultDBDSN = "postgres://user:password@db:5432/note_tracker?sslmode=disable"
	// defaultMigratePath = "migrations"
)

func ReadConfig() (*Config, error) {
	var cfg Config

	flag.StringVar(&cfg.Host, "host", defaultHost, "flag for explicit server host specifications")
	flag.IntVar(&cfg.Port, "port", defaultPort, "flag for explicit server port specifications")
	flag.BoolVar(&cfg.Debug, "debug", false, "flag for explicit debug mode")
	flag.StringVar(&cfg.DBDSN, "db", defaultDBDSN, "flag for explicit database specifications")
	// flag.StringVar(&cfg.MigratePath, "migrate", defaultMigratePath, "flag for explicit migrate path specifications")

	flag.Parse()

	if cfg.Host == "localhost" {
		cfg.Host = cmp.Or(os.Getenv("HOST"), cfg.Host)
	}

	if cfg.Port == defaultPort {
		defPort := strconv.Itoa(cfg.Port)
		envPort := cmp.Or(os.Getenv("PORT"), defPort)
		port, err := strconv.Atoi(envPort)
		if err != nil {
			return nil, err
		}
		cfg.Port = port
	}

	cfg.DBDSN = cmp.Or(os.Getenv("DB_DSN"), cfg.DBDSN)

	// if cfg.MigratePath == defaultMigratePath {
	// 	cfg.MigratePath = cmp.Or(os.Getenv("MIGRATE_PATH"), cfg.MigratePath)
	// }

	return &cfg, nil

}
