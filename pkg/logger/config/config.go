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
	ProjectName string `env:"PROJECT_NAME,default=golang_project" yaml:"projectName"`
	// log library name
	Code string `env:"LOG_CODE,default=logrus" yaml:"logCode"`
	// log encoding
	Encoding string `env:"LOG_ENCODING,default=console" yaml:"logEncoding"`
	// logging level
	LogLevel string `env:"LOG_LEVEL,default=info" yaml:"logLevel"`
	// date time format
	DateTimeFormat string `env:"DATE_TIME_FORMAT,default:2006-01-02 15:04:05" yaml:"dateTimeFormat"`
	// date format
	DateFormat string `env:"DATE_FORMAT,default:2006-01-02" yaml:"dateFormat"`
	// show caller in log message
	EnableCaller bool `yaml:"enableCaller"`
	// development mode marker
	DevMode bool `env:"DEV_MODE,default=false" yaml:"devMode"`
	// output writer
	Out io.Writer
}
