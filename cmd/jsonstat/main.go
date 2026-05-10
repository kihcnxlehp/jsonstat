package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kihcnxhelp/jsonstat/internal/config"
	"github.com/kihcnxhelp/jsonstat/internal/processor"
)

func run() error {
	cfg, err := config.Load()
	if errors.Is(err, config.ErrHelpRequested) {
		flag.Usage()
		os.Exit(0)
	}
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	log.Printf("processing %s (max=%d filter=%s=%s)",
		cfg.InputFile, cfg.MaxRecords, cfg.FilterField, cfg.FilterValue)

	stats, err := processor.Process(cfg.InputFile, cfg.FilterField, cfg.FilterValue, cfg.MaxRecords)
	if err != nil {
		return fmt.Errorf("process: %w", err)
	}

	fmt.Printf("total=%d matched=%d skipped=%d\n",
		stats.Total, stats.Matched, stats.Skipped)

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
