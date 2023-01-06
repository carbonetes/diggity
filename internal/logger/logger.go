package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

// GetLogger returns a logger
func GetLogger() *logrus.Logger {
	return &logrus.Logger{
		Out: os.Stderr,
		ExitFunc: func(i int) {
			os.Exit(i)
		},
		Level: logrus.DebugLevel,
		Formatter: &easy.Formatter{
			LogFormat: "%msg%\n",
		},
	}
}

// SetQuietMode sets logger to quiet mode
func SetQuietMode(logger *logrus.Logger) *logrus.Logger {
	logger.SetOutput(io.Discard)
	return logger
}
