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
type Transport string

const (
	DEVELOPMENT AppMode = "DEVELOPMENT"
	PRODUCTION  AppMode = "PRODUCTION"

	TRANSPORT_HTTP Transport = "HTTP"
	TRANSPORT_GRPC Transport = "GRPC"
)

type Config struct {
	Transport Transport
	Logging   logConfig.Logging `yaml:"logging"`
	Project   ProjectConfig     `yaml:"project"`
	MongoDB   MongoDBConfig     `yaml:"mongodb"`
	Http      HttpConfig        `yaml:"http"`
	Storage   StorageConfig     `yaml:"storage"`
}

type ProjectConfig struct {
	Name                   string `env:"PROJECT_NAME" yaml:"name"`
	Mode                   string `env:"APPLICATION_MODE"`
	Version                string `env:"APPLICATION_VERSION" yaml:"version"`
	Salt                   string `env:"APP_SALT"`
	GracefulTimeoutSeconds int    `yaml:"gracefulTimeoutSeconds"`
	JwtSecret              string `env:"APPLICATION_JWT_SECRET"`
	SwaggerEnabled         bool   `yaml:"swaggerEnabled"`
	FileUploadMaxMegabytes int    `yaml:"fileUploadMaxMegabytes"`
}

type MongoDBConfig struct {
	URI      string `env:"MONGODB_URI"`
	User     string `env:"MONGODB_USER"`
	Password string `env:"MONGODB_PASSWORD"`
	Database string `yaml:"database"`
}

type HttpConfig struct {
	Port                  int
	AccessTokenTTLMinutes int `yaml:"accessTokenTTLMinutes"`
	RefreshTokenTTLHours  int `yaml:"refreshTokenTTLHours"`
}

type StorageConfig struct {
	URI             string `env:"STORAGE_URI"`
	User            string `env:"STORAGE_MINIO_USER"`
	Password        string `env:"STORAGE_MINIO_PASSWORD"`
	AccessKeyID     string `env:"STORAGE_ACCESS_KEY_ID"`
	SecretAccessKey string `env:"STORAGE_SECRET_ACCESS_KEY"`
	BucketName      string `yaml:"bucketName"`
}

func Load(transport Transport) *Config {
	cfg := Config{Transport: transport}

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
