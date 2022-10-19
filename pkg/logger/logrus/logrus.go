// Package logrus handles creating logrus logger
package logrus

import (
	"log"
	"os"

	"github.com/abdukhashimov/student_aggregator/pkg/logger"
	"github.com/abdukhashimov/student_aggregator/pkg/logger/config"

	"github.com/sirupsen/logrus"
)

// logWrapper allows to add or override logger methods
type logWrapper struct {
	*logrus.Logger
}

// RegisterLogrusLog creates a new logrus logger
func RegisterLogrusLog(cfg *config.Logging) (logger.Logger, error) {
	logrusLog := logrus.New()
	logrusLog.SetFormatter(&logrus.TextFormatter{})
	logrusLog.SetReportCaller(true)
	log.SetOutput(os.Stdout)
	err := customizeLogrusLogFromConfig(logrusLog, cfg)
	if err != nil {
		return nil, err
	}

	wrapper := logWrapper{
		logrusLog,
	}

	return &wrapper, nil
}

// customizeLogrusLogFromConfig customizes log based on parameters from configuration
func customizeLogrusLogFromConfig(log *logrus.Logger, cfg *config.Logging) error {
	log.SetReportCaller(cfg.EnableCaller)
	l := &log.Level
	err := l.UnmarshalText([]byte(cfg.LogLevel))
	if err != nil {
		return err
	}
	log.SetLevel(*l)

	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = cfg.DateTimeFormat
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)

	if cfg.Out != nil {
		log.SetOutput(cfg.Out)
	}

	return nil
}

// Sync method is required for logger.Logger interface. logrus doesn't have Sync method
func (_ *logWrapper) Sync() error {
	return nil
}
