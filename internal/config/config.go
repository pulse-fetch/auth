package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type GrpcConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}
type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	DBName   string `yaml:"db_name"`
	Password string `yaml:"password"`
	Sslmode  string `yaml:"sslmode"`
}
type Config struct {
	Env         string         `yaml:"env"`
	HmacSecret  string         `yaml:"hmac-secret"`
	GrpcCfg     GrpcConfig     `yaml:"grpc-server"`
	PostgresCfg PostgresConfig `yaml:"postgres"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("empty config path")
	}
	cfg, err := LoadWithPath(path)
	if err != nil {
		panic("Failed parsing config, error: " + err.Error())
	}
	return cfg
}

func LoadWithPath(path string) (*Config, error) {
	var cfg Config
	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func fetchConfigPath() string {
	var path string
	flag.StringVar(&path, "config", "", "this is config path")
	flag.Parse()
	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}
	return path
}
