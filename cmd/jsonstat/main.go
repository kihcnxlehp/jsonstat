package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/kihcnxhelp/jsonstat/internal/config"
	"github.com/kihcnxhelp/jsonstat/internal/processor"
)

func run() error {
	cfg, err := config.Load()
	if errors.Is(err, config.ErrHelpRequested) {
		flag.Usage()
		return nil
	}
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	log.Printf("processing %s (max=%d filter=%s=%s)",
		cfg.InputFile, cfg.MaxRecords, cfg.FilterField, cfg.FilterValue)

	var input io.Reader

	if cfg.InputFile == "-" {
		input = os.Stdin
	} else {
		file, err := os.Open(cfg.InputFile)
		if err != nil {
			return fmt.Errorf("open input file: %w", err)
		}

		defer file.Close()

		input = file
	}

	stats, err := processor.Process(input, os.Stdout, cfg.FilterField, cfg.FilterValue, cfg.MaxRecords)
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
