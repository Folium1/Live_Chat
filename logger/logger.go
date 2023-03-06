package logger

import (
	"fmt"
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	// create a new logrus logger
	Log = logrus.New()

	// set the log output to stdout
	Log.SetOutput(os.Stdout)

	// set the log level to debug
	Log.SetLevel(logrus.DebugLevel)

	// set the formatter
	Log.SetFormatter(&logrus.TextFormatter{})
}

// log an error message
func Error(msg string, err error, funcName string) {
	Log.WithError(err).Error(msg, " | function = "+funcName)
}

// log an informational message
func Info(msg string, funcName string) {
	Log.Info(msg, " | function = "+funcName)
}

func InfoHttp(path, method, funcName string) {
	msg := fmt.Sprintf("%v | %v | %v", path, method, funcName)
	Log.Info(msg)
}

func GetFuncName() string {
	pc, _, _, _ := runtime.Caller(2)
	function := runtime.FuncForPC(pc).Name()
	return function
}
