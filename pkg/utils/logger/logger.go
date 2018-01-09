package logger

import (
	"runtime"

	log "github.com/Sirupsen/logrus"
)

// DecorateLoggerWithRuntimeContext appends line, file and function context to the logger
func DecorateLoggerWithRuntimeContext(logger *log.Entry) *log.Entry {
	if pc, file, line, ok := runtime.Caller(1); ok {
		fName := runtime.FuncForPC(pc).Name()
		return logger.WithField("file", file).WithField("line", line).WithField("func", fName)
	}
	return logger
}
