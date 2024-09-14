package logging

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
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
	level := strings.ToUpper(entry.Level.String())
	timestamp := entry.Time.Format(time.RFC3339Nano)
	colorFunc := color.New(setOutputColorPerLevel(level)).SprintFunc()
	logMessage := fmt.Sprintf("[%s] %s: %s\n", colorFunc(level), timestamp, entry.Message)
	return []byte(logMessage), nil
}

func setOutputColorPerLevel(level string) color.Attribute {
	var selectedColor color.Attribute
	switch level {
	case "DEBUG":
		selectedColor = color.FgCyan
	case "INFO":
		selectedColor = color.FgGreen
	case "WARN", "WARNING":
		selectedColor = color.FgYellow
	case "ERROR", "PANIC", "FATAL":
		selectedColor = color.FgRed
	default:
		selectedColor = color.FgWhite
	}
	return selectedColor
}
