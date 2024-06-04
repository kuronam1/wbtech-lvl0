package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type Config struct {
	Env        string `yaml:"env"`
	HttpServer `yaml:"httpServer"`
	PG         `yaml:"pg"`
	NatsStream `yaml:"natsStream"`
}

type HttpServer struct {
	Host            string        `yaml:"host"`
	Port            string        `yaml:"port"`
	Timeout         time.Duration `yaml:"timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdownTimeout"`
}

type PG struct {
	Login    string `yaml:"login"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DBName   string `yaml:"dbName"`
	SslMode  string `yaml:"sslMode"`
	MaxConn  int    `yaml:"maxConn"`
	MinConn  int    `yaml:"minConn"`
}

type NatsStream struct {
	ClusterID string `yaml:"clusterID"`
	ClientID  string `yaml:"clientID"`
}

func NewConfig(path string) (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func New(configPath string) (*Config, error) {
	cfg := &Config{}
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
