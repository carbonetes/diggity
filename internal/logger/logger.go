package logger

import (
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

var (
	log     = logrus.New()
	logFile = "/logs/app.log"
)

func init() {
	if !CheckFileIfExists(logFile) {
		err := os.MkdirAll(filepath.Dir(logFile), 0700)
		if err != nil {
			log.Fatal(err)
		}
	}

	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(io.MultiWriter(file, os.Stderr))

	log.SetLevel(logrus.DebugLevel)
	log.SetReportCaller(true)
	log.SetFormatter(&Formatter{})
}

func GetLogger() *logrus.Logger {
	return log
}

func CheckFileIfExists(file string) bool {
	if _, err := os.Stat(file); err == nil {
		return true
	}
	return false
}
