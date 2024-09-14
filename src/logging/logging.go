package logging

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var stringToLogLevel map[string]logrus.Level

func init() {
	stringToLogLevel = map[string]logrus.Level{
		"DEBUG": logrus.DebugLevel,
		"INFO":  logrus.InfoLevel,
		"WARN":  logrus.WarnLevel,
		"ERROR": logrus.ErrorLevel,
		"PANIC": logrus.PanicLevel,
		"FATAL": logrus.FatalLevel,
		"TRACE": logrus.TraceLevel,
	}
}

func ConfigureLogger(level string) {
	// Set logging level
	loglevel, ok := stringToLogLevel[strings.ToUpper(level)]
	if ok {
		logrus.SetLevel(loglevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	logrus.SetReportCaller(true)
	logrus.SetFormatter(&CustomFormatter{})
}

type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Create a custom format for the log message
	logMessage := fmt.Sprintf("[%s] %s: %s\n", entry.Level.String(), entry.Time.Format(time.RFC3339Nano), entry.Message)
	return []byte(logMessage), nil
}
