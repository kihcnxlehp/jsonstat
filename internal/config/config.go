package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
)

var ErrHelpRequested = errors.New("help requested")

var validFields = map[string]struct{}{
	"id":     {},
	"name":   {},
	"role":   {},
	"salary": {},
}

var validLevels = map[string]struct{}{
	"debug": {},
	"info":  {},
	"warn":  {},
	"error": {},
}

// Config содержит всю конфигурацию утилиты.
type Config struct {
	InputFile   string
	FilterField string
	FilterValue string
	MaxRecords  int
	LogLevel    string
}

type fileConfig struct {
	LogLevel   string `json:"log_level"`
	MaxRecords *int   `json:"max_records"`
}

// loadFile загружает конфиг
func loadFile(path string) (fileConfig, error) {
	var fileCfg fileConfig

	data, err := os.ReadFile(path)
	if err != nil {
		return fileCfg, fmt.Errorf("failed to read config file: %w", err)
	}

	if len(data) == 0 {
		return fileCfg, nil
	}

	err = json.Unmarshal(data, &fileCfg)
	if err != nil {
		return fileCfg, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return fileCfg, nil
}

// Validate валидирует конфиг
func (c Config) Validate() error {
	if c.MaxRecords < 0 {
		return fmt.Errorf("max records must be greater than or equal to 0: %d", c.MaxRecords)
	}

	if _, ok := validLevels[c.LogLevel]; !ok {
		return fmt.Errorf("invalid log level: %s", c.LogLevel)
	}

	if (c.FilterField != "" && c.FilterValue == "") || (c.FilterField == "" && c.FilterValue != "") {
		return fmt.Errorf("both -field and -value must be specified together")
	}

	if c.FilterField != "" {
		if _, ok := validFields[c.FilterField]; !ok {
			return fmt.Errorf("invalid filter field: %s", c.FilterField)
		}
	}

	return nil
}

// Load строит конфигурацию по цепочке: defaults -> file -> env -> flags.
func Load() (Config, error) {
	// 1. Defaults
	cfg := Config{
		InputFile:   "input.json",
		FilterField: "",
		FilterValue: "",
		MaxRecords:  0,
		LogLevel:    "info",
	}

	// 2. File
	configPath := os.Getenv("JSONSTAT_CONFIG")
	if configPath == "" {
		configPath = "config.json"
	}

	fileCfg, err := loadFile(configPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return cfg, fmt.Errorf("load config: %w", err)
		}
	} else {
		if fileCfg.LogLevel != "" {
			cfg.LogLevel = fileCfg.LogLevel
		}
		if fileCfg.MaxRecords != nil {
			cfg.MaxRecords = *fileCfg.MaxRecords
		}
	}

	// 3. Env
	if v := os.Getenv("JSONSTAT_INPUT_FILE"); v != "" {
		cfg.InputFile = v
	}
	if v := os.Getenv("JSONSTAT_FILTER_FIELD"); v != "" {
		cfg.FilterField = v
	}
	if v := os.Getenv("JSONSTAT_FILTER_VALUE"); v != "" {
		cfg.FilterValue = v
	}
	if v := os.Getenv("JSONSTAT_MAX_RECORDS"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return cfg, fmt.Errorf("invalid JSONSTAT_MAX_RECORDS: %w", err)
		}

		cfg.MaxRecords = n
	}
	if v := os.Getenv("JSONSTAT_LOG_LEVEL"); v != "" {
		cfg.LogLevel = v
	}

	// 4. Flags
	flag.StringVar(&cfg.InputFile, "input", cfg.InputFile, "input JSON Lines file")
	flag.StringVar(&cfg.FilterField, "field", cfg.FilterField, "field to filter by")
	flag.StringVar(&cfg.FilterValue, "value", cfg.FilterValue, "value to filter for")
	flag.IntVar(&cfg.MaxRecords, "max", cfg.MaxRecords, "max records to process (0 = unlimited)")
	flag.StringVar(&cfg.LogLevel, "log-level", cfg.LogLevel, "log level (debug/info/warn/error)")

	// Help flag
	help := flag.Bool("help", false, "show help")

	flag.Parse()

	if *help {
		return cfg, ErrHelpRequested
	}

	return cfg, nil
}
