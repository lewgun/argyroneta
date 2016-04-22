package logger

import (
	"os"

	"github.com/Sirupsen/logrus"
)

func New(format, level string) *logrus.Logger {
	log := logrus.New()
	log.Out = os.Stdout
	log.Formatter = newFormatter(format)
	log.Level = parseLogLevel(level)
	return log
}

func newFormatter(format string) logrus.Formatter {
	switch format {
	case "text", "":
		return &logrus.TextFormatter{}
	case "json":
		return &logrus.JSONFormatter{}
	default:
		return &logrus.TextFormatter{}
	}
}

func parseLogLevel(level string) logrus.Level {
	switch level {
	case "panic":
		return logrus.PanicLevel
	case "fatal":
		return logrus.FatalLevel
	case "error":
		return logrus.ErrorLevel
	case "warn", "warning":
		return logrus.WarnLevel
	case "info":
		return logrus.InfoLevel
	case "debug":
		return logrus.DebugLevel
	default:
		return logrus.InfoLevel
	}
}
