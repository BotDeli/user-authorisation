package config

import (
	"fmt"
	"time"
)

type Config struct {
	Grpc     *GRPCConfig     `yaml:"grpc" env-required:"true"`
	Postgres *PostgresConfig `yaml:"postgres" env-required:"true"`
	Redis    *RedisConfig    `yaml:"redis" env-required:"true"`
}

type GRPCConfig struct {
	Network string `yaml:"network" env-default:"tcp"`
	Address string `yaml:"address" env-required:"true"`
}

type PostgresConfig struct {
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password"`
	Address  string `yaml:"address" env-required:"true"`
	Dbname   string `yaml:"dbname" env-required:"true"`
	Sslmode  string `yaml:"sslmode" env-default:"false"`
}

func (cfg *PostgresConfig) GetSourceName() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Address,
		cfg.Dbname,
		cfg.Sslmode,
	)
	//return "host=127.0.0.1 port=5432 user=kare dbname=postgres sslmode=require"
}

type RedisConfig struct {
	Network       string        `yaml:"network" env-default:"tcp"`
	Address       string        `yaml:"address" env-default:"localhost:6379"`
	Password      string        `yaml:"password" env-default:""`
	DB            int           `yaml:"db" env-default:"0"`
	LifetimeWrite time.Duration `yaml:"lifetimeWrite" env-default:"24h"`
}
