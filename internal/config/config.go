package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env"`
	Logger     `yaml:"logger"`
	HTTPServer `yaml:"http_server"`
	DB         `yaml:"db"`
	Auth       `yaml:"auth"`
	SecretKey  []byte
}

type Auth struct {
	AccessTokenTTL  time.Duration `yaml:"access_token_ttl"`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl"`
	Issuer          string        `yaml:"issuer"`
}

type Logger struct {
	Level string `yaml:"level"`
}

const (
	EnvProd  = "prod"
	EnvDev   = "dev"
	EnvLocal = "local"
)

const (
	LogLevelDebug  = "debug"
	LogLevelInfo   = "info"
	LogLevelWarn   = "warn"
	LogLevelError  = "error"
	LogLevelDPanic = "dpanic"
	LogLevelPanic  = "panic"
	LogLevelFatal  = "fatal"
)

type HTTPServer struct {
	Address     string        `yaml:"address"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

type DB struct {
	Host            string        `yaml:"host"`
	Port            string        `yaml:"port"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	Name            string        `yaml:"name"`
	SSLMode         string        `yaml:"ssl_mode"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time"`
}

func MustLoadConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("CONFIG_PATH does not exist: %s", configPath)
	}

	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatal("Cannot read config: &v", err)
	}

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("SECRET_KEY environment variable not set")
	}

	config.SecretKey = []byte(secretKey)

	return &config
}
