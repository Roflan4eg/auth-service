package config

import (
	"fmt"
	"time"
)

type DatabaseConfig struct {
	Host     string `yaml:"-" env:"HOST"`
	Port     string `yaml:"-" env:"PORT"`
	User     string `yaml:"-" env:"USER"`
	Password string `yaml:"-" env:"PASSWORD"`
	Name     string `yaml:"-" env:"NAME"`
	SSLMode  string `yaml:"ssl_mode" env:"SSLMODE" envDefault:"disable"`

	MaxConns        int           `yaml:"max_conns" env:"MAX_CONNS" envDefault:"10"`
	MaxConnLifetime time.Duration `yaml:"max_conn_lifetime" env:"MAX_CONN_LIFETIME" envDefault:"1h"`
	ConnectTimeout  time.Duration `yaml:"connect_timeout" env:"CONNECT_TIMEOUT" envDefault:"5s"`
	//MinConns          int         `yaml:"min_conns" json:"min_conns" env:"DB_MIN_CONNS"`
	//MaxConnIdleTime   time.Duration `yaml:"max_conn_idle_time" json:"max_conn_idle_time" env:"DB_MAX_CONN_IDLE_TIME"`
	//HealthCheckPeriod time.Duration `yaml:"health_check_period" json:"health_check_period" env:"DB_HEALTH_CHECK_PERIOD"`
	//StatementTimeout time.Duration `yaml:"statement_timeout" json:"statement_timeout" env:"DB_STATEMENT_TIMEOUT"`
	//QueryTimeout     time.Duration `yaml:"query_timeout" json:"query_timeout" env:"DB_QUERY_TIMEOUT"`
}

func (c *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode)
}

type HTTPConfig struct {
	Port         string        `yaml:"port" env:"PORT" envDefault:"8080"`
	Host         string        `yaml:"host" env:"HOST" envDefault:"0.0.0.0"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env:"READ_TIMEOUT" envDefault:"30s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env:"WRITE_TIMEOUT" envDefault:"30s"`
}

func (c *HTTPConfig) Address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

type GRPCConfig struct {
	Port         string        `yaml:"port" env:"PORT" envDefault:"9090"`
	Host         string        `yaml:"host" env:"HOST" envDefault:"0.0.0.0"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env:"READ_TIMEOUT" envDefault:"30s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env:"WRITE_TIMEOUT" envDefault:"30s"`
}

func (c *GRPCConfig) Address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

type RedisConfig struct {
	Host     string        `yaml:"-" env:"HOST" envDefault:"localhost"`
	Port     string        `yaml:"-" env:"PORT" envDefault:"6379"`
	Password string        `yaml:"-" env:"PASSWORD" envDefault:""`
	DB       int           `yaml:"-" env:"DB" envDefault:"0"`
	TTL      time.Duration `yaml:"ttl" env:"TTL" envDefault:"168h"`
}

func (c *RedisConfig) Address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

type JWTConfig struct {
	Secret          string        `yaml:"-" env:"SECRET_KEY"`
	AccessTokenTTL  time.Duration `yaml:"access_token_ttl" env:"ACCESS_TOKEN_TTL" envDefault:"15m"`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl" env:"REFRESH_TOKEN_TTL" envDefault:"1h"`
}

type AppConfig struct {
	Name            string        `yaml:"name" env:"NAME" envDefault:"auth-service"`
	Environment     string        `yaml:"environment" env:"ENV" envDefault:"local"`
	Version         string        `yaml:"version" env:"VERSION" envDefault:"0.0.1"`
	LogPath         string        `yaml:"log_path" env:"LOG_PATH" envDefault:"stdout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env:"SHUTDOWN_TIMEOUT" envDefault:"10s"`
	AutoMigrate     bool          `yaml:"auto_migrate" env:"AUTO_MIGRATE" envDefault:"true"`
}

type Config struct {
	App       *AppConfig      `yaml:"app" envPrefix:"APP_"`
	Database  *DatabaseConfig `yaml:"database" envPrefix:"DB_"`
	HTTP      *HTTPConfig     `yaml:"http" envPrefix:"HTTP_"`
	Redis     *RedisConfig    `yaml:"redis" envPrefix:"REDIS_"`
	GRPC      *GRPCConfig     `yaml:"grpc" envPrefix:"GRPC_"`
	JWTConfig *JWTConfig      `yaml:"jwt" envPrefix:"JWT_"`
}
