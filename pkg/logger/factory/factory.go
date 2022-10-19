package factory

import (
	"fmt"

	"github.com/abdukhashimov/student_aggregator/pkg/logger"
	"github.com/abdukhashimov/student_aggregator/pkg/logger/config"
	"github.com/abdukhashimov/student_aggregator/pkg/logger/logrus"
	"github.com/abdukhashimov/student_aggregator/pkg/logger/zap"
)

// logger map to map logger code to logger builder
var logFactoryBuilderMap = map[string]loggerBuilder{
	config.LOGRUS: &logrus.Factory{},
	config.ZAP:    &zap.Factory{},
}

// interface for logger factory
type loggerBuilder interface {
	Build(cfg *config.Logging) (logger.Logger, error)
}

// accessors for factoryBuilderMap
func getLogFactoryBuilder(key string) (loggerBuilder, error) {
	logFactoryBuilder, ok := logFactoryBuilderMap[key]
	if !ok {
		return nil, fmt.Errorf("not supported logger: %s", key)
	}

	return logFactoryBuilder, nil
}

// Build logger using appropriate log factory
func Build(cfg *config.Logging) (logger.Logger, error) {
	logFactoryBuilder, err := getLogFactoryBuilder(cfg.Code)
	if err != nil {
		return nil, err
	}

	log, err := logFactoryBuilder.Build(cfg)
	if err != nil {
		return nil, err
	}

	return log, nil
}
