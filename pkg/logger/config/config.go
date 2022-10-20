package config

import "io"

// constant for logger code, it needs to match log code (logConfig)in configuration
const (
	LOGRUS string = "logrus"
	ZAP    string = "zap"
)

const (
	DEBUG = "debug"
	INFO  = "info"
	WARN  = "warn"
	ERROR = "error"
)

// Logging represents logger handler
type Logging struct {
	// project name
	ProjectName string `yaml:"projectName"`
	// log library name
	Code string `yaml:"code"`
	// log encoding
	Encoding string `yaml:"encoding"`
	// logging level
	LogLevel string `yaml:"level"`
	// date time format
	DateTimeFormat string `yaml:"dateTimeFormat"`
	// date format
	DateFormat string `yaml:"dateFormat"`
	// show caller in log message
	EnableCaller bool `yaml:"enableCaller"`
	// development mode marker
	DevMode bool `yaml:"devMode"`
	// output writer
	Out io.Writer
}
