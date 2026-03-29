package config

import "fmt"

func Validate(cfg Config) error {
	if cfg.Service.Name == "" {
		return fmt.Errorf("SERVICE_NAME is required")
	}
	if cfg.HTTP.Addr == "" {
		return fmt.Errorf("HTTP_ADDR is required")
	}
	if cfg.DB.DSN == "" {
		return fmt.Errorf("DB_DSN is required")
	}
	if cfg.DB.MaxConns < 1 {
		return fmt.Errorf("DB_MAX_CONNS must be >= 1")
	}
	return nil
}
