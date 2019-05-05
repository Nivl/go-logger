// Package logger contains interfaces to deal with loggers
package logger

// gomock interface, requires mockgen
// Update with "go generate github.com/Nivl/go-logger"
//go:generate mockgen -destination mocklogger/logger.go -package mocklogger github.com/Nivl/go-logger Logger

// Logger is an interface used for all loggers
type Logger interface {
	// ID returns the logger's unique ID
	ID() string

	// Close frees any resource allocated by the logger
	// the logger may not be reusable after being closed
	Close() error

	// IsClosed returns wether the logger is closed or not
	IsClosed() bool

	// Error logs an error message
	// Arguments are handled in the manner of fmt.Println.
	Error(msg string)

	// Info logs a message that may be helpful, but isn’t essential,
	// for troubleshooting
	// Arguments are handled in the manner of fmt.Println.
	Info(msg string)

	// Debug logs a message that is intended for use in a development
	// environment while actively debugging your subsystem, not in shipping
	// software
	// Arguments are handled in the manner of fmt.Println.
	Debug(msg string)

	// Log logs a message that might result a failure
	// Arguments are handled in the manner of fmt.Println.
	Log(msg string)
}
