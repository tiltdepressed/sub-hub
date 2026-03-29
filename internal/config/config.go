package config

import "time"

type Config struct {
	Env     string        `env:"ENV" env-default:"local"`
	Service ServiceConfig `env-prefix:"SERVICE_"`
	Log     LogConfig     `env-prefix:"LOG_"`
	HTTP    HTTPConfig    `env-prefix:"HTTP_"`
	DB      DBConfig      `env-prefix:"DB_"`
	Tracing TracingConfig `env-prefix:"TRACING_"`
	OpenAPI OpenAPIConfig `env-prefix:"OPENAPI_"`
}

type ServiceConfig struct {
	Name string `env:"NAME" env-default:"sub-hub"`
}

type LogConfig struct {
	Level string `env:"LEVEL" env-default:"info"`
}

type HTTPConfig struct {
	Addr         string        `env:"ADDR" env-default:":8080"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT" env-default:"5s"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" env-default:"10s"`
	IdleTimeout  time.Duration `env:"IDLE_TIMEOUT" env-default:"60s"`
}

type DBConfig struct {
	DSN             string        `env:"DSN" env-default:""`
	ConnectTimeout  time.Duration `env:"CONNECT_TIMEOUT" env-default:"5s"`
	MaxConns        int32         `env:"MAX_CONNS" env-default:"10"`
	MinConns        int32         `env:"MIN_CONNS" env-default:"0"`
	HealthcheckPing time.Duration `env:"HEALTHCHECK_PING" env-default:"1s"`
}

type TracingConfig struct {
	Enabled     bool   `env:"ENABLED" env-default:"false"`
	Endpoint    string `env:"ENDPOINT" env-default:"localhost:4317"`
	ServiceName string `env:"SERVICE_NAME" env-default:"sub-hub"`
}

type OpenAPIConfig struct {
	Path string `env:"PATH" env-default:"api/openapi.yaml"`
}
