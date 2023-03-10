package config

import (
	"errors"
	"log"
	"time"

	"github.com/spf13/viper"
)

var (
	cfg *Config
)

type Config struct {
	Server         ServerConfig
	Postgres       PostgresConfig
	Jwt            JwtConfig
	FirstSuperUser FirstSuperUserConfig
	Logger         Logger
}

// Server config struct
type ServerConfig struct {
	AppVersion   string
	Port         string
	Mode         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// Logger config
type Logger struct {
	Encoding string
	Level    string
}

// Postgresql config
type PostgresConfig struct {
	PostgresqlHost     string
	PostgresqlPort     string
	PostgresqlUser     string
	PostgresqlPassword string
	PostgresqlDbname   string
	PostgresqlSSLMode  string
}

// Jwt config
type JwtConfig struct {
	// SecretKey                string
	Issuer                        string
	JwtAccessTokenExpireDuration  int64
	JwtAccessTokenPrivateKey      string
	JwtAccessTokenPublicKey       string
	JwtRefreshTokenExpireDuration int64
	JwtRefreshTokenPrivateKey     string
	JwtRefreshTokenPublicKey      string
}

// First Super User
type FirstSuperUserConfig struct {
	Email    string
	Name     string
	Password string
}

// Load config file from given path
func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

// Parse config file
func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	cfg = &c

	return &c, nil
}

func GetCfg() *Config {
	if cfg == nil {
		cfg = new(Config)
	}
	return cfg
}
