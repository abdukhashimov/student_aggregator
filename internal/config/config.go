package config

import (
	"fmt"
	"os"

	logConfig "github.com/abdukhashimov/student_aggregator/pkg/logger/config"

	env "github.com/Netflix/go-env"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type AppMode string

const (
	DEVELOPMENT AppMode = "DEVELOPMENT"
	PRODUCTION  AppMode = "PRODUCTION"
)

type Config struct {
	Logging logConfig.Logging `yaml:"logging"`
	Project struct {
		Name                   string `env:"PROJECT_NAME" yaml:"name"`
		Mode                   string `env:"APPLICATION_MODE"`
		Version                string `env:"APPLICATION_VERSION" yaml:"version"`
		Salt                   string `env:"APP_SALT"`
		GracefulTimeoutSeconds int    `yaml:"gracefulTimeoutSeconds"`
	} `yaml:"project"`
	MongoDB struct {
		URI      string `env:"MONGODB_URI"`
		User     string `env:"MONGODB_USER"`
		Password string `env:"MONGODB_PASSWORD"`
		Database string `yaml:"database"`
	} `yaml:"mongodb"`
	Http struct {
		Port int
	}
}

func Load() *Config {
	var cfg Config

	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	appMode := getAppMode()
	configPath, err := getConfigPath(appMode)
	if err != nil {
		panic(err)
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		panic(err)
	}

	_, err = env.UnmarshalFromEnviron(&cfg)
	if err != nil {
		panic("unmarshal from environment error")
	}

	if err := validateConfig(&cfg); err != nil {
		panic(err)
	}

	return &cfg
}

func getAppMode() AppMode {
	mode := AppMode(os.Getenv("APPLICATION_MODE"))
	if mode != PRODUCTION {
		mode = DEVELOPMENT
	}

	return mode
}

func getConfigPath(appMode AppMode) (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}
	suffix := "Dev"
	if appMode == PRODUCTION {
		suffix = "Prod"
	}

	return fmt.Sprintf("%s/configs/appConfig%s.yaml", path, suffix), nil
}
