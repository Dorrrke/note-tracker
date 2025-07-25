// Package config предоставляет функционал для загрузки конфигурации приложения
// из различных источников: файлов, флагов командной строки и переменных окружения.
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
)

// Config описывает структуру конфигурации приложения.
type Config struct {
	Host  string `json:"host"`
	Port  int    `json:"port"`
	Debug bool   `json:"debug"`
	DBDSN string `json:"dbdsn"`
}

const (
	defaultPort  = 8080
	defaultHost  = "localhost"
	defaultDBDSN = "postgres://user:password@db:5432/note_tracker?sslmode=disable"
)

// getConfigFromFile загружает конфигурацию из JSON-файла по указанному пути.
// Возвращает указатель на Config или ошибку, если файл не удалось открыть или распарсить.
func getConfigFromFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable open config file: %w", err)
	}
	defer file.Close()

	var cfg Config
	if err = json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("unable decode data in config file: %w", err)
	}

	return &cfg, nil
}

// parseFlags регистрирует и парсит флаги командной строки.
// Значения сохраняются в переданную структуру cfg и строку configFile.
func parseFlags(cfg *Config, configFile *string) {
	flag.StringVar(&cfg.Host, "host", defaultHost, "flag for explicit server host specifications")
	flag.IntVar(&cfg.Port, "port", defaultPort, "flag for explicit server port specifications")
	flag.BoolVar(&cfg.Debug, "debug", false, "flag for explicit debug mode")
	flag.StringVar(&cfg.DBDSN, "db", defaultDBDSN, "flag for explicit database specifications")

	flag.StringVar(configFile, "c", "", "path to config file (JSON format)")
	flag.StringVar(configFile, "config", "", "path to config file (JSON format)")

	flag.Parse()
}

// applyConfigFile применяет параметры из файла конфигурации,
// если они не были переопределены через флаги.
func applyConfigFile(cfg *Config, path string) error {
	if path == "" {
		return nil
	}

	fileCfg, err := getConfigFromFile(path)
	if err != nil {
		return err
	}

	if cfg.Host == defaultHost {
		cfg.Host = fileCfg.Host
	}
	if cfg.Port == defaultPort {
		cfg.Port = fileCfg.Port
	}
	if !cfg.Debug {
		cfg.Debug = fileCfg.Debug
	}
	if cfg.DBDSN == defaultDBDSN {
		cfg.DBDSN = fileCfg.DBDSN
	}

	return nil
}

// applyEnvVars применяет значения из переменных окружения,
// если они заданы.
func applyEnvVars(cfg *Config) error {
	if host := os.Getenv("HOST"); host != "" {
		cfg.Host = host
	}
	if portStr := os.Getenv("PORT"); portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return fmt.Errorf("invalid PORT value: %w", err)
		}
		cfg.Port = port
	}
	if debugStr := os.Getenv("DEBUG"); debugStr != "" {
		debug, err := strconv.ParseBool(debugStr)
		if err != nil {
			return fmt.Errorf("invalid DEBUG value: %w", err)
		}
		cfg.Debug = debug
	}
	if dbdsn := os.Getenv("DB_DSN"); dbdsn != "" {
		cfg.DBDSN = dbdsn
	}
	return nil
}

// ReadConfig читает конфигурацию приложения, объединяя данные из флагов,
// файла конфигурации и переменных окружения. Приоритет источников:
// 1. Переменные окружения
// 2. Файл конфигурации
// 3. Флаги командной строки
// 4. Значения по умолчанию
//
// Возвращает указатель на структуру Config или ошибку при её формировании.
func ReadConfig() (*Config, error) {
	cfg := &Config{
		Host:  defaultHost,
		Port:  defaultPort,
		Debug: false,
		DBDSN: defaultDBDSN,
	}

	var configFile string
	parseFlags(cfg, &configFile)
	if err := applyConfigFile(cfg, configFile); err != nil {
		return nil, err
	}
	if err := applyEnvVars(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
