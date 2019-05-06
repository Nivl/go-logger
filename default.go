package logger

var defaultManager = NewManager()

// AddGlobalData is used to add data that will be added to all logs
func AddGlobalData(key string, value interface{}) {
	defaultManager.AddGlobalData(key, value)
}

// RemoveGlobalData is used to remove data that are added to all logs
func RemoveGlobalData(key string) {
	defaultManager.RemoveGlobalData(key)
}

// Add adds a logger
// returns ErrAlreadyExist if the logger has already been added
func Add(l Logger) error {
	return defaultManager.Add(l)
}

// Remove safely removes a logger
// returns an Error struct if a logger could not be safely removed.
// Upon errors the logger will be force removed from the manager
func Remove(loggerID string) error {
	return defaultManager.Remove(loggerID)
}

// Close safely removes all the loggers
// All submanagers will also be closed
// returns a list of errors if a logger could not be safely removed.
func Close() []error {
	return defaultManager.Close()
}

// NewSubManager creates a new manager that can have its own loggers.
// The tag of the current manager will be passed to the submanager.
// Calling a logging method on a submanager will trigger the same logging
// method on the parent.
func NewSubManager(tag string) Manager {
	return defaultManager.NewSubManager(tag)
}

// SetTag adds a tag to the logs
func SetTag(tag string) {
	defaultManager.SetTag(tag)
}

// Tag returns the tag of the manager
func Tag() string {
	return defaultManager.Tag()
}

// FullTag returns the full tag (including parents) of the manager
func FullTag() string {
	return defaultManager.FullTag()
}

// ID returns the manager's unique ID
func ID() string {
	return defaultManager.ID()
}

// Errorf logs an error message
// Arguments are handled in the manner of fmt.Printf
func Errorf(msg string, args ...interface{}) {
	defaultManager.Errorf(msg, args...)
}

// Error logs an error message
// Arguments are handled in the manner of fmt.Println.
func Error(args ...interface{}) {
	defaultManager.Error(args...)
}

// Infof logs a message that may be helpful, but isn’t essential,
// for troubleshooting
// Arguments are handled in the manner of fmt.Printf
func Infof(msg string, args ...interface{}) {
	defaultManager.Infof(msg, args...)
}

// Info logs a message that may be helpful, but isn’t essential,
// for troubleshooting
// Arguments are handled in the manner of fmt.Println.
func Info(args ...interface{}) {
	defaultManager.Info(args...)
}

// Debugf logs a message that is intended for use in a development
// environment while actively debugging your subsystem, not in shipping
// software
// Arguments are handled in the manner of fmt.Printf
func Debugf(msg string, args ...interface{}) {
	defaultManager.Debugf(msg, args...)
}

// Debug logs a message that is intended for use in a development
// environment while actively debugging your subsystem, not in shipping
// software
// Arguments are handled in the manner of fmt.Println.
func Debug(args ...interface{}) {
	defaultManager.Debug(args...)
}

// Logf logs a message that might result a failure
// Arguments are handled in the manner of fmt.Printf
func Logf(msg string, args ...interface{}) {
	defaultManager.Logf(msg, args...)
}

// Log logs a message that might result a failure
// Arguments are handled in the manner of fmt.Println.
func Log(args ...interface{}) {
	defaultManager.Log(args...)
}
