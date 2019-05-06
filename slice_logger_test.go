package logger

// we make sure SliceLogger implements Logger
var _ Logger = (*SliceLogger)(nil)

// NewSliceLogger creates and returns a slice logger
func NewSliceLogger() Logger {
	return &SliceLogger{}
}

// SliceLogger is a logger that puts everything in a slice (useful for testing)
// /!\ Not go-routine-safe
type SliceLogger struct {
	data   []string
	closed bool
	id     string
}

func (l *SliceLogger) ID() string {
	if l.id != "" {
		return l.id
	}
	return "slice-logger"
}

func (l *SliceLogger) Close() error {
	l.data = []string{}
	l.closed = true
	return nil
}

func (l *SliceLogger) clear() {
	l.data = []string{}
}

func (l *SliceLogger) IsClosed() bool {
	return l.closed
}

func (l *SliceLogger) Error(msg string) {
	l.write(msg, levelError)
}

func (l *SliceLogger) Info(msg string) {
	l.write(msg, levelInfo)
}

func (l *SliceLogger) Debug(msg string) {
	l.write(msg, levelDebug)
}

func (l *SliceLogger) Log(msg string) {
	l.write(msg, levelDefault)
}

func (l *SliceLogger) write(msg string, lvl logLevel) {
	msg = lvl.Tag() + msg
	l.data = append(l.data, msg)
}
