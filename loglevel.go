package logger

import (
	"fmt"
)

// logLevel is an helper type to hold all the different log levels
type logLevel int

// ALl the log levels
const (
	levelDefault logLevel = iota
	levelDebug
	levelInfo
	levelError
)

func (level logLevel) Tag() string {
	levelStr := ""

	switch level {
	case levelDebug:
		levelStr = "DEBUG"
	case levelInfo:
		levelStr = "INFO"
	case levelError:
		levelStr = "ERROR"
	default:
		return ""
	}

	return fmt.Sprintf("[%s]", levelStr)
}
