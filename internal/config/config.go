package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	HttpServer `yaml:"httpServer" env-required:"true"`
	PG         `yaml:"pg" env-required:"true"`
	NatsStream `yaml:"natsStream" env-required:"true"`
}

type HttpServer struct {
	Host            string        `yaml:"host" env-default:"localhost"`
	Port            string        `yaml:"port" env-default:":8080"`
	Timeout         time.Duration `yaml:"timeout" env-default:"4s"`
	ShutdownTimeout time.Duration `yaml:"shutdownTimeout" env-default:"10s"`
}

type PG struct {
	Login    string `yaml:"login"`
	Password string `yaml:"password"`
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:":5432"`
	DBName   string `yaml:"DBName"`
	SslMode  string `yaml:"sslMode" env-default:"disabled"`
	MaxConn  int    `yaml:"maxConn" env-default:"20"`
	MinConn  int    `yaml:"minConn" env-default:"5"`
}

type NatsStream struct {
	ClusterID string
	ClientID  string
	NatsUrl   string
}

func NewConfig(path string) (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
