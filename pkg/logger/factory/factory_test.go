package factory

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/abdukhashimov/student_aggregator/pkg/logger"
	"github.com/abdukhashimov/student_aggregator/pkg/logger/config"
)

const DateTimeFormat = "2006-01-02 15:04:05"

const (
	DEBUG  = "debug"
	INFO   = "info"
	WARN   = "warn"
	ERROR  = "error"
	DEBUGF = "debugf"
	INFOF  = "infof"
	WARNF  = "warnf"
	ERRORF = "errorf"
)

var logrusConfig = config.Logging{
	Code:           config.LOGRUS,
	LogLevel:       config.DEBUG,
	DateTimeFormat: DateTimeFormat,
	EnableCaller:   false,
	Out:            nil,
}

var zapConfig = config.Logging{
	ProjectName:    "Student Aggregator",
	Code:           config.ZAP,
	LogLevel:       config.DEBUG,
	DateTimeFormat: "2006-01-02 15:04:05",
	DateFormat:     "2006-01-02",
	Encoding:       "json",
	DevMode:        true,
}

type TestCase struct {
	level            string
	executionMethod  func(args ...interface{})
	executionMethodF func(msg string, args ...interface{})
	expected         bool
}

type ZapLog struct {
	L string
	M string
}

var oldStdout *os.File

func buildLogger(lc config.Logging) (*os.File, *os.File, logger.Logger, error) {
	oldStdout = os.Stdout
	reader, writer, _ := os.Pipe()
	os.Stdout = writer
	lc.Out = os.Stdout
	log, err := Build(&lc)
	if err != nil {
		return nil, nil, nil, err
	}

	return reader, writer, log, nil
}

func callLoger(reader *os.File, writer *os.File, cb func()) string {
	outC := make(chan string)
	go func() {
		scanner := bufio.NewScanner(reader)
		if scanner.Scan() {
			outC <- scanner.Text()
		} else {
			outC <- ""
		}
	}()

	go func() {
		cb()
		os.Stdout = oldStdout
		writer.Close()
	}()

	result := <-outC
	close(outC)

	return result
}

func buildTestCases(level string) map[string]func(log logger.Logger) TestCase {
	switch level {
	case DEBUG:
		return map[string]func(log logger.Logger) TestCase{
			DEBUG:  func(log logger.Logger) TestCase { return TestCase{DEBUG, log.Debug, nil, true} },
			DEBUGF: func(log logger.Logger) TestCase { return TestCase{DEBUG, nil, log.Debugf, true} },
			INFO:   func(log logger.Logger) TestCase { return TestCase{INFO, log.Info, nil, true} },
			INFOF:  func(log logger.Logger) TestCase { return TestCase{INFO, nil, log.Infof, true} },
			WARN:   func(log logger.Logger) TestCase { return TestCase{WARN, log.Warn, nil, true} },
			WARNF:  func(log logger.Logger) TestCase { return TestCase{WARN, nil, log.Warnf, true} },
			ERROR:  func(log logger.Logger) TestCase { return TestCase{ERROR, log.Error, nil, true} },
			ERRORF: func(log logger.Logger) TestCase { return TestCase{ERROR, nil, log.Errorf, true} },
		}
	case INFO:
		return map[string]func(log logger.Logger) TestCase{
			DEBUG:  func(log logger.Logger) TestCase { return TestCase{DEBUG, log.Debug, nil, false} },
			DEBUGF: func(log logger.Logger) TestCase { return TestCase{DEBUG, nil, log.Debugf, false} },
			INFO:   func(log logger.Logger) TestCase { return TestCase{INFO, log.Info, nil, true} },
			INFOF:  func(log logger.Logger) TestCase { return TestCase{INFO, nil, log.Infof, true} },
			WARN:   func(log logger.Logger) TestCase { return TestCase{WARN, log.Warn, nil, true} },
			WARNF:  func(log logger.Logger) TestCase { return TestCase{WARN, nil, log.Warnf, true} },
			ERROR:  func(log logger.Logger) TestCase { return TestCase{ERROR, log.Error, nil, true} },
			ERRORF: func(log logger.Logger) TestCase { return TestCase{ERROR, nil, log.Errorf, true} },
		}
	case WARN:
		return map[string]func(log logger.Logger) TestCase{
			DEBUG:  func(log logger.Logger) TestCase { return TestCase{DEBUG, log.Debug, nil, false} },
			DEBUGF: func(log logger.Logger) TestCase { return TestCase{DEBUG, nil, log.Debugf, false} },
			INFO:   func(log logger.Logger) TestCase { return TestCase{INFO, log.Info, nil, false} },
			INFOF:  func(log logger.Logger) TestCase { return TestCase{INFO, nil, log.Infof, false} },
			WARN:   func(log logger.Logger) TestCase { return TestCase{WARN, log.Warn, nil, true} },
			WARNF:  func(log logger.Logger) TestCase { return TestCase{WARN, nil, log.Warnf, true} },
			ERROR:  func(log logger.Logger) TestCase { return TestCase{ERROR, log.Error, nil, true} },
			ERRORF: func(log logger.Logger) TestCase { return TestCase{ERROR, nil, log.Errorf, true} },
		}
	case ERROR:
		return map[string]func(log logger.Logger) TestCase{
			DEBUG:  func(log logger.Logger) TestCase { return TestCase{DEBUG, log.Debug, nil, false} },
			DEBUGF: func(log logger.Logger) TestCase { return TestCase{DEBUG, nil, log.Debugf, false} },
			INFO:   func(log logger.Logger) TestCase { return TestCase{INFO, log.Info, nil, false} },
			INFOF:  func(log logger.Logger) TestCase { return TestCase{INFO, nil, log.Infof, false} },
			WARN:   func(log logger.Logger) TestCase { return TestCase{WARN, log.Warn, nil, false} },
			WARNF:  func(log logger.Logger) TestCase { return TestCase{WARN, nil, log.Warnf, false} },
			ERROR:  func(log logger.Logger) TestCase { return TestCase{ERROR, log.Error, nil, true} },
			ERRORF: func(log logger.Logger) TestCase { return TestCase{ERROR, nil, log.Errorf, true} },
		}
	}

	return nil
}

func TestLogrusLogger(t *testing.T) {
	for _, level := range []string{DEBUG, INFO, WARN, ERROR} {
		t.Run(fmt.Sprintf("%s level test", level), func(t *testing.T) {
			testCaseLogrusConfig := logrusConfig
			testCaseLogrusConfig.LogLevel = level
			testCasesMap := buildTestCases(level)
			if testCasesMap == nil {
				t.Errorf("unsopported level")
				return
			}

			for method, caseBuilder := range testCasesMap {
				reader, writer, log, err := buildLogger(testCaseLogrusConfig)
				if err != nil {
					t.Errorf("unexpected error")
					return
				}
				testCase := caseBuilder(log)
				msg := fmt.Sprintf("example %s", method)
				result := callLoger(reader, writer, func() {
					if testCase.executionMethod != nil {
						testCase.executionMethod(msg)
					} else {
						msg = "example %s"
						testCase.executionMethodF(msg, method)
					}
				})

				if testCase.expected {
					if !strings.Contains(result, fmt.Sprintf("level=%s", testCase.level)) {
						t.Errorf("log doesn't containst appropriate log level. logger level - %s, method - %s, method level - %s", level, testCase.level, method)
					}

					msg = fmt.Sprintf("example %s", method)
					if !strings.Contains(result, fmt.Sprintf(`msg="%s"`, msg)) {
						t.Errorf("log doesn't containst log message. logger level - %s, method - %s, method level - %s", level, testCase.level, method)
					}
				} else {
					if result != "" {
						t.Errorf("not expected method calling. logger level - %s, method - %s, method level - %s", level, testCase.level, method)
					}
				}

			}
		})

	}
}

func TestZapLogger(t *testing.T) {
	for _, level := range []string{DEBUG, INFO, WARN, ERROR} {
		t.Run(fmt.Sprintf("%s level test", level), func(t *testing.T) {
			testCaseZapConfig := zapConfig
			testCaseZapConfig.LogLevel = level
			testCasesMap := buildTestCases(level)
			if testCasesMap == nil {
				t.Errorf("unsopported level")
				return
			}
			for method, caseBuilder := range testCasesMap {
				reader, writer, log, err := buildLogger(testCaseZapConfig)
				if err != nil {
					os.Stdout = oldStdout
					writer.Close()
					t.Errorf("unexpected error")
					return
				}
				testCase := caseBuilder(log)
				msg := fmt.Sprintf("example %s", method)
				result := callLoger(reader, writer, func() {
					if testCase.executionMethod != nil {
						testCase.executionMethod(msg)
					} else {
						msg = "example %s"
						testCase.executionMethodF(msg, method)
					}
				})
				if testCase.expected {
					var parsedResult ZapLog
					if err := json.Unmarshal([]byte(result), &parsedResult); err != nil {
						t.Errorf("expected valid JSON. logger level - %s, method - %s, method level - %s", level, testCase.level, method)
					}

					if parsedResult.L != strings.ToUpper(testCase.level) {
						t.Errorf("log doesn't containst appropriate log level. logger level - %s, method - %s, method level - %s", level, testCase.level, method)
					}

					msg = fmt.Sprintf("example %s", method)
					if parsedResult.M != msg {
						t.Errorf("log doesn't containst log message. logger level - %s, method - %s, method level - %s", level, testCase.level, method)
					}
				} else {
					if result != "" {
						t.Errorf("not expected method calling. logger level - %s, method - %s, method level - %s", level, testCase.level, method)
					}
				}
			}
		})
	}
}
