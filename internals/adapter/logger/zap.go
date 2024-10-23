package logger

import (
	"github.com/jocbarbosa/viswals-backend/internals/core/port"
	"go.uber.org/zap"
)

type zapAdapter struct {
	logger  *zap.SugaredLogger
	err     error
	details []port.Detail
}

// NewZapAdapter creates a new zap adapter for the Logger interface
func NewZapAdapter(logger *zap.Logger) port.Logger {
	return &zapAdapter{
		logger: logger.Sugar(),
	}
}

// Debug logs a debug message
func (z *zapAdapter) Debug(msg string, args ...interface{}) {
	z.logger.Debugw(msg, z.withDetails(args...)...)
}

// Info logs an info message
func (z *zapAdapter) Info(msg string, args ...interface{}) {
	z.logger.Infow(msg, z.withDetails(args...)...)
}

// Error logs an error message
func (z *zapAdapter) Error(msg string, args ...interface{}) {
	z.logger.Errorw(msg, z.withDetails(args...)...)
}

// Alarm logs an alarm message
func (z *zapAdapter) Alarm(msg string, args ...interface{}) {
	z.logger.Errorw("ALARM: "+msg, z.withDetails(args...)...)
}

// WithError adds an error to the logger
func (z *zapAdapter) WithError(err error) port.Logger {
	newLogger := *z
	newLogger.err = err
	return &newLogger
}

// WithDetails adds details to the logger
func (z *zapAdapter) WithDetails(details ...port.Detail) port.Logger {
	newLogger := *z
	newLogger.details = append(newLogger.details, details...)
	return &newLogger
}

// withDetails adds error and details to the logger
func (z *zapAdapter) withDetails(args ...interface{}) []interface{} {
	fields := make([]interface{}, 0)
	if z.err != nil {
		fields = append(fields, "error", z.err)
	}
	for _, detail := range z.details {
		fields = append(fields, detail.Key, detail.Value)
	}
	fields = append(fields, args...)
	return fields
}
