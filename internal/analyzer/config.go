package analyzer

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	CheckLowercase     bool     `yaml:"check_lowercase"`
	CheckEnglish       bool     `yaml:"check_english"`
	CheckSpecialChars  bool     `yaml:"check_special_chars"`
	CheckSensitiveData bool     `yaml:"check_sensitive_data"`
	SensitivePatterns  []string `yaml:"sensitive_patterns"`
	AllowedLoggerNames []string `yaml:"allowed_logger_names"`
}

func defaultConfig() Config {
	return Config{
		CheckLowercase:     true,
		CheckEnglish:       true,
		CheckSpecialChars:  true,
		CheckSensitiveData: true,
		SensitivePatterns: []string{
			"password",
			"passwd",
			"token",
			"secret",
			"api_key",
			"apikey",
			"access_key",
		},
		AllowedLoggerNames: []string{
			"slog",
			"logger",
		},
	}
}

func loadConfig(path string) Config {
	cfg := defaultConfig()

	if path == "" {
		path = ".loglint.yml"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return defaultConfig()
	}

	if len(cfg.SensitivePatterns) == 0 {
		cfg.SensitivePatterns = defaultConfig().SensitivePatterns
	}

	if len(cfg.AllowedLoggerNames) == 0 {
		cfg.AllowedLoggerNames = defaultConfig().AllowedLoggerNames
	}

	return cfg
}
