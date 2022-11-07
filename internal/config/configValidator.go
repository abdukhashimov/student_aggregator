package config

import (
	"fmt"
)

func validateConfig(cfg *Config) error {
	if err := validateGeneralConfig(cfg); err != nil {
		return err
	}

	if err := validateMongoDbConfig(cfg); err != nil {
		return err
	}

	if err := validateLoggingConfig(cfg); err != nil {
		return err
	}

	if err := validateStorageConfig(cfg); err != nil {
		return err
	}

	return nil
}

func validateGeneralConfig(cfg *Config) error {
	if cfg.Project.Name == "" {
		return buildError("project name")
	}

	if cfg.Project.Mode == "" {
		return buildError("project mode")
	}

	if cfg.Project.Version == "" {
		return buildError("project version")
	}

	if cfg.Project.Salt == "" {
		return buildError("project salt")
	}

	if cfg.Transport == TRANSPORT_HTTP {
		if cfg.Project.JwtSecret == "" {
			return buildError("jtw secret")
		}

		if cfg.Http.AccessTokenTTLMinutes == 0 {
			return buildError("access token ttl")
		}

		if cfg.Http.RefreshTokenTTLHours == 0 {
			return buildError("refresh token ttl")
		}
	}

	return nil
}

func validateMongoDbConfig(cfg *Config) error {
	if cfg.MongoDB.URI == "" {
		return buildError("mongodb uri")
	}

	if cfg.MongoDB.Database == "" {
		return buildError("database")
	}

	return nil
}

func validateLoggingConfig(cfg *Config) error {
	if cfg.Logging.Code == "" {
		return buildError("logger type")
	}

	if cfg.Logging.LogLevel == "" {
		return buildError("logger level")
	}

	if cfg.Logging.DateTimeFormat == "" {
		return buildError("logger datetime format")
	}

	return nil
}

func validateStorageConfig(cfg *Config) error {
	if cfg.Storage.URI == "" {
		return buildError("storage user format")
	}

	if cfg.Storage.AccessKeyID == "" {
		return buildError("storage access key ID format")
	}
	if cfg.Storage.SecretAccessKey == "" {
		return buildError("storage secret access key format")
	}

	if cfg.Storage.BucketName == "" {
		return buildError("storage bucket name format")
	}

	if cfg.Project.Mode == string(DEVELOPMENT) {
		if cfg.Storage.User == "" {
			return buildError("storage user format")
		}
		if len(cfg.Storage.Password) < 8 {
			return buildError("storage password min length is 8")
		}
	}

	return nil
}

func buildError(key string) error {
	return fmt.Errorf("application config. %s is not specified", key)
}
