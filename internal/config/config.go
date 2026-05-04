package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
)

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
	MaxRecords int    `json:"max_records"`
}

// loadFile загружает кофиг
func loadFile(path string) (fileConfig, error) {
	var fileCfg fileConfig

	file, err := os.ReadFile(path)
	if err != nil {
		return fileCfg, fmt.Errorf("failed to read config file: %w", err)
	}

	err = json.Unmarshal(file, &fileCfg)
	if err != nil {
		fmt.Println("Error unmarshalling file: ", err)
		return fileCfg, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return fileCfg, nil
}

// Load строит конфигурацию по цепочке: defaults -> file -> env -> flags.
func Load() Config {
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
		fmt.Println("Error loading config: ", err)
	} else {
		if fileCfg.LogLevel != "" {
			cfg.LogLevel = fileCfg.LogLevel
		}
		if fileCfg.MaxRecords != 0 {
			cfg.MaxRecords = fileCfg.MaxRecords
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
		if n, convertErr := strconv.Atoi(v); convertErr == nil {
			cfg.MaxRecords = n
		}
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

	flag.Visit(func(f *flag.Flag) {
	})

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	return cfg
}
