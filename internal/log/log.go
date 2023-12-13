package log

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

var log = logrus.New()

func init() {
	// Send all logs to nowhere by default
	log.SetOutput(io.Discard)

	// Send all logs with level higher than warning to stderr
	log.AddHook(&writer.Hook{
		Writer: os.Stderr,
		LogLevels: []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		},
	})

	// Send info and debug logs to stdout
	log.AddHook(&writer.Hook{
		Writer: os.Stdout,
		LogLevels: []logrus.Level{
			logrus.InfoLevel,
			logrus.DebugLevel,
		},
	})

	// Set the exit function to exit with code 0
	log.ExitFunc = func(_ int) {
		os.Exit(0)
	}

	log.SetFormatter(&easy.Formatter{
		LogFormat: "%msg%\n",
	})
}

// Print func prints the arguments to stdout
func Print(arg ...interface{}) {
	log.Print(arg...)
}

// Printf func prints the formatted arguments to stdout
func Printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

// Error func prints the arguments to stderr1
func Error(arg ...interface{}) {
	log.Error(arg...)
}

// Errorf func prints the formatted arguments to stderr
func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

// Info func prints the arguments to stdout
func Info(arg ...interface{}) {
	log.Info(arg...)
}

// Infof func prints the formatted arguments to stdout
func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

// Debug func prints the arguments to stdout
func Debug(arg ...interface{}) {
	log.Debug(arg...)
}

// Debugf func prints the formatted arguments to stdout
func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

// Fatal func prints the arguments to stderr and exits with code 0
func Fatal(arg ...interface{}) {
	log.Fatal(arg...)
}

// Fatalf func prints the formatted arguments to stderr and exits with code 0
func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}
