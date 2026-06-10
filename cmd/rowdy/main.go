package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/webdevmeg42/rowdy/pkg/validator/config"
	"github.com/webdevmeg42/rowdy/pkg/validator/db"
	"github.com/webdevmeg42/rowdy/pkg/validator/reporter"
	"github.com/webdevmeg42/rowdy/pkg/validator/runner"
)

func main() {
	configPath := flag.String("config", "", "path to YAML config file (required)")
	format := flag.String("format", "terminal", "output format: terminal or junit")
	dbOverride := flag.String("db", "", "override database.path from config")
	flag.Parse()

	if *configPath == "" {
		fmt.Fprintln(os.Stderr, "error: --config is required")
		flag.Usage()
		os.Exit(2)
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}
	if *dbOverride != "" {
		cfg.Database.Path = *dbOverride
	}

	conn, err := db.Open(cfg.Database.Path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}
	defer conn.Close()

	results := runner.Run(conn, cfg.TestCases)

	switch *format {
	case "junit":
		if err := reporter.JUnit(os.Stdout, results); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(2)
		}
	default:
		reporter.Terminal(os.Stdout, results)
	}

	for _, r := range results {
		if !r.Passed {
			os.Exit(1)
		}
	}
}
