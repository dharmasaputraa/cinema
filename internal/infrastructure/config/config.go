package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App   AppConfig
	DB    DBConfig
	Redis RedisConfig
	JWT   JWTConfig
}

type AppConfig struct {
	Env         string
	Port        string
	CORSOrigins []string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret             string
	ExpiryHours        int
	RefreshExpiryHours int
}

func Load() (*Config, error) {
	env := getEnv("APP_ENV", "development")

	if env == "production" {
		viper.SetConfigFile(".env.prod")
	} else {
		viper.SetConfigFile(".env.dev")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("No .env file found, using system environment variables")
	}

	cfg := &Config{
		App: AppConfig{
			Env:         viper.GetString("APP_ENV"),
			Port:        viper.GetString("APP_PORT"),
			CORSOrigins: parseCSV(viper.GetString("CORS_ORIGINS")),
		},
		DB: DBConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSLMODE"),
		},
		Redis: RedisConfig{
			Host:     viper.GetString("REDIS_HOST"),
			Port:     viper.GetString("REDIS_PORT"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
		JWT: JWTConfig{
			Secret:             viper.GetString("JWT_SECRET"),
			ExpiryHours:        viper.GetInt("JWT_EXPIRY_HOURS"),
			RefreshExpiryHours: viper.GetInt("JWT_REFRESH_EXPIRY_HOURS"),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func parseCSV(value string) []string {
	if value == "" {
		return []string{}
	}
	parts := strings.Split(value, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func (c *Config) Validate() error {
	var errs []string

	// ── APP ─────────────────────────────────────
	if c.App.Env == "" {
		errs = append(errs, "APP_ENV is required")
	}

	if c.App.Port == "" {
		errs = append(errs, "APP_PORT is required")
	}

	// ── DATABASE ───────────────────────────────
	if c.DB.Host == "" {
		errs = append(errs, "DB_HOST is required")
	}
	if c.DB.Port == "" {
		errs = append(errs, "DB_PORT is required")
	}
	if c.DB.User == "" {
		errs = append(errs, "DB_USER is required")
	}
	if c.DB.Name == "" {
		errs = append(errs, "DB_NAME is required")
	}

	// ── REDIS ──────────────────────────────────
	if c.Redis.Host == "" {
		errs = append(errs, "REDIS_HOST is required")
	}
	if c.Redis.Port == "" {
		errs = append(errs, "REDIS_PORT is required")
	}

	// ── JWT ────────────────────────────────────
	if c.JWT.Secret == "" {
		errs = append(errs, "JWT_SECRET is required")
	}
	if c.JWT.ExpiryHours <= 0 {
		errs = append(errs, "JWT_EXPIRY_HOURS must be > 0")
	}
	if c.JWT.RefreshExpiryHours <= 0 {
		errs = append(errs, "JWT_REFRESH_EXPIRY_HOURS must be > 0")
	}

	// ── PRODUCTION STRICT RULES  ─────────────
	if c.App.Env == "production" {
		if len(c.App.CORSOrigins) == 0 {
			errs = append(errs, "CORS_ORIGINS must be set in production")
		}

		if c.DB.Password == "" {
			errs = append(errs, "DB_PASSWORD is required in production")
		}

		if len(c.JWT.Secret) < 16 {
			errs = append(errs, "JWT_SECRET must be at least 16 characters in production")
		}
	}

	// ── FINAL CHECK ────────────────────────────
	if len(errs) > 0 {
		return fmt.Errorf("config validation error:\n- %s", strings.Join(errs, "\n- "))
	}

	return nil
}
