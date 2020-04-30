// This file and its contents are licensed under the Apache License 2.0.
// Please see the included NOTICE for copyright information and
// LICENSE for a copy of the license

// Package log creates logs in the same way as Prometheus, while ignoring errors
package log

import (
	"fmt"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/common/promlog"
)

var (
	// Application wide logger
	logger log.Logger
)

// Init starts logging given the minimum log level
func Init(logLevel string) error {
	allowedLevel := promlog.AllowedLevel{}
	err := allowedLevel.Set(logLevel)
	if err != nil {
		return err
	}

	config := promlog.Config{
		Level:  &allowedLevel,
		Format: &promlog.AllowedFormat{},
	}

	logger = promlog.New(&config)
	return nil
}

// Debug logs a DEBUG level message, ignoring logging errors
func Debug(keyvals ...interface{}) {
	_ = level.Debug(logger).Log(keyvals...)
}

// Info logs an INFO level message, ignoring logging errors
func Info(keyvals ...interface{}) {
	_ = level.Info(logger).Log(keyvals...)
}

// Warn logs a WARN level message, ignoring logging errors
func Warn(keyvals ...interface{}) {
	_ = level.Warn(logger).Log(keyvals...)
}

// Error logs an ERROR level message, ignoring logging errors
func Error(keyvals ...interface{}) {
	_ = level.Error(logger).Log(keyvals...)
}

// Fatal logs an ERROR level message and exits
func Fatal(keyvals ...interface{}) {
	_ = level.Error(logger).Log(keyvals...)
	os.Exit(1)
}

// CustomCacheLogger is a custom logger used for transforming cache logs
// so that they conform the our logging setup. It also implements the
// bigcache.Logger interface.
type CustomCacheLogger struct{}

// Printf sends the log in the debug stream of the logger.
func (c *CustomCacheLogger) Printf(format string, v ...interface{}) {
	_ = level.Debug(logger).Log("msg", fmt.Sprintf(format, v...))
}
