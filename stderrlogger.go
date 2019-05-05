package logger

import (
	"log"
)

// we make sure StderrLogger implements Logger
var _ Logger = (*StderrLogger)(nil)

// NewStderrLogger creates and returns a logger that writes on stderr
func NewStderrLogger() Logger {
	return &StderrLogger{}
}

// StderrLogger is a non-buffered logger that writes on stderr
type StderrLogger struct {
	closed bool
}

// ID returns the logger's unique ID
func (l *StderrLogger) ID() string {
	return "stderr-logger"
}

// Close frees any resource allocated by the logger
// the logger may not be reusable after being closed
func (l *StderrLogger) Close() error {
	l.closed = true
	return nil
}

// IsClosed returns wether the logger is closed or not
func (l *StderrLogger) IsClosed() bool {
	return l.closed
}

// Error logs an error message
// Arguments are handled in the manner of fmt.Println.
func (l *StderrLogger) Error(msg string) {
	l.write(msg, levelError)
}

// Info logs a message that may be helpful, but isn’t essential,
// for troubleshooting
// Arguments are handled in the manner of fmt.Println.
func (l *StderrLogger) Info(msg string) {
	l.write(msg, levelInfo)
}

// Debug logs a message that is intended for use in a development
// environment while actively debugging your subsystem, not in shipping
// software
// Arguments are handled in the manner of fmt.Println.
func (l *StderrLogger) Debug(msg string) {
	l.write(msg, levelDebug)
}

// Log logs a message that might result a failure
// Arguments are handled in the manner of fmt.Println.
func (l *StderrLogger) Log(msg string) {
	l.write(msg, levelDefault)
}

func (l *StderrLogger) write(msg string, lvl logLevel) {
	msg = lvl.Tag() + msg
	log.Print(msg)
}
