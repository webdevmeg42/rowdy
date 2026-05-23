package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database  DatabaseConfig `yaml:"database"`
	TestCases []TestCase     `yaml:"test_cases"`
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

type TestCase struct {
	Name       string      `yaml:"name"`
	Seed       *SeedConfig `yaml:"seed,omitempty"`
	Query      string      `yaml:"query"`
	Assertions []Assertion `yaml:"assertions"`
}

type SeedConfig struct {
	Table string                   `yaml:"table"`
	Rows  []map[string]interface{} `yaml:"rows"`
}

type Assertion struct {
	Type     string      `yaml:"type"`
	Expected *int        `yaml:"expected,omitempty"`
	Column   string      `yaml:"column,omitempty"`
	Row      int         `yaml:"row,omitempty"`
	Equals   interface{} `yaml:"equals,omitempty"`
	Format   string      `yaml:"format,omitempty"`
}

var validTypes = map[string]bool{
	"row_count": true, "column_exists": true, "not_null": true,
	"value": true, "format": true,
}

var validFormats = map[string]bool{
	"email": true, "uuid": true, "date": true, "url": true,
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	if err := validate(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func validate(cfg *Config) error {
	if cfg.Database.Path == "" {
		return fmt.Errorf("database.path is required")
	}
	for i, tc := range cfg.TestCases {
		if tc.Name == "" {
			return fmt.Errorf("test_cases[%d]: name is required", i)
		}
		if tc.Query == "" {
			return fmt.Errorf("test_cases[%d] %q: query is required", i, tc.Name)
		}
		for j, a := range tc.Assertions {
			if !validTypes[a.Type] {
				return fmt.Errorf("test_cases[%d] %q: assertions[%d]: unknown type %q", i, tc.Name, j, a.Type)
			}
			if a.Type == "format" && !validFormats[a.Format] {
				return fmt.Errorf("test_cases[%d] %q: assertions[%d]: unknown format %q", i, tc.Name, j, a.Format)
			}
			if a.Type == "row_count" && a.Expected == nil {
				return fmt.Errorf("test_cases[%d] %q: assertions[%d]: row_count requires expected", i, tc.Name, j)
			}
			if (a.Type == "column_exists" || a.Type == "not_null") && a.Column == "" {
				return fmt.Errorf("test_cases[%d] %q: assertions[%d]: %s requires column", i, tc.Name, j, a.Type)
			}
			if a.Type == "value" && a.Column == "" {
				return fmt.Errorf("test_cases[%d] %q: assertions[%d]: value requires column", i, tc.Name, j)
			}
			if a.Type == "value" && a.Equals == nil {
				return fmt.Errorf("test_cases[%d] %q: assertions[%d]: value requires equals", i, tc.Name, j)
			}
		}
	}
	return nil
}
