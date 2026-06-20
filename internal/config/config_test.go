package config

import "testing"

func TestConfig_Validate(t *testing.T) {
	valid := Config{
		InputFile:  "input.json",
		MaxRecords: 10,
		LogLevel:   "info",
	}

	tests := []struct {
		name    string
		modify  func(*Config)
		wantErr bool
	}{
		{"valid config", func(c *Config) {}, false},
		{"negative max records", func(c *Config) { c.MaxRecords = -1 }, true},
		{"invalid log level", func(c *Config) { c.LogLevel = "invalid" }, true},
		{"field without value", func(c *Config) { c.FilterField = "role" }, true},
		{"value without field", func(c *Config) { c.FilterValue = "admin" }, true},
		{"both value and field - ok", func(c *Config) { c.FilterField = "role"; c.FilterValue = "admin" }, false},
		{"invalid filter field", func(c *Config) { c.FilterField = "email"; c.FilterValue = "x" }, true},
		{"stdin marker is valid", func(c *Config) { c.InputFile = "-" }, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := valid
			tt.modify(&cfg)

			err := cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
