package mocks

import (
	"github.com/abdukhashimov/student_aggregator/internal/pkg/logger"
)

type mockLogger struct {
}

func (m mockLogger) Debugf(format string, args ...interface{}) {
}

func (m mockLogger) Infof(format string, args ...interface{}) {
}

func (m mockLogger) Warnf(format string, args ...interface{}) {
}

func (m mockLogger) Errorf(format string, args ...interface{}) {
}

func (m mockLogger) Fatalf(format string, args ...interface{}) {
}

func (m mockLogger) Debug(args ...interface{}) {
}

func (m mockLogger) Info(args ...interface{}) {
}

func (m mockLogger) Warn(args ...interface{}) {
}

func (m mockLogger) Error(args ...interface{}) {
}

func (m mockLogger) Fatal(args ...interface{}) {
}

func (m mockLogger) Sync() error {
	return nil
}

func MockAppLogger() {
	log := &mockLogger{}
	logger.SetLogger(log)
}
