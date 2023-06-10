package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Database struct {
	Username string `yaml:"postgres.username"`
	Password string `yaml:"postgres.password"`
	Host     string `yaml:"postgres.host"`
	Port     string `yaml:"postgres.port"`
	DB       string `yaml:"postgres.db"`
	SslMode  string `yaml:"postgres.sslmode"`
	TimeZone string `yaml:"postgres.timezone"`
}

type Service struct {
	Name    string `yaml:"service.name"`
	ID      uint32 `yaml:"service.id"`
	BaseURL string `yaml:"service.baseURL"`
	HTTP    struct {
		Host           string `yaml:"http.host"`
		Port           string `yaml:"http.port"`
		RequestTimeout uint   `yaml:"http.requestTimeout"`
		AllowOrigin    string `yaml:"http.allowOrigin"`
	}
	Token struct {
		Password   string `yaml:"token.password"`
		Expiration uint   `yaml:"token.expiration"`
	}
}

type Config struct {
	Debug    bool     // if true we run on debug mode
	Service  Service  `yaml:"service"`
	Postgres Database `yaml:"database"`
}

type ConfigInterface interface {
	Set(key string, query []byte) error
	SetDebug(bool)
	GetDebug() bool
	Load(path string) error
}

var Confs = Config{}

// Set method
// you can set new key in switch for manage config with config server
func (g *Config) Set(key string, query []byte) error {
	if err := json.Unmarshal(query, &Confs); err != nil {
		return err
	}
	return nil
}

func (g *Config) GetDebug() bool {
	return g.Debug
}

func (g *Config) SetDebug(debug bool) {
	g.Debug = debug
}

// Load returns configs
func (g *Config) Load(path string) error {
	if path == "" {
		dir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("Can't get working directory with error: %e", err)
		}
		path = dir + "/config/config-debug.yaml"
		if mode := os.Getenv("GIN_MODE"); mode == "release" {
			path = dir + "/config/config.yaml"
			g.Debug = false
		}
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return g.file(path)
	}

	return fmt.Errorf("file not exists")
}

// file func
func (g *Config) file(path string) error {

	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return err

	}

	return viper.Unmarshal(&Confs)
}
