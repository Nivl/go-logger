package logger

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// gomock interface, requires mockgen
// Update with "go generate github.com/Nivl/go-logger"
//go:generate mockgen -destination mocklogger/manager.go -package mocklogger github.com/Nivl/go-logger Manager

// List of all errors
var (
	ErrAlreadyExist = errors.New("logger already added")
)

// Manager is an interface used to manage loggers
type Manager interface {
	// ID returns the manager's unique ID
	ID() string

	// AddGlobalData is used to add data that will be added to all logs
	AddGlobalData(key string, value interface{})

	// RemoveGlobalData is used to remove data that are added to all logs
	RemoveGlobalData(key string)

	// Add adds a new logger
	// returns ErrAlreadyExist if the logger has already been added
	Add(Logger) error

	// Remove safely removes a logger
	// returns the logger and an error if the logger could not be safely remove.
	// Upon errors the logger will be force removed from the manager
	Remove(loggerID string) error

	// Close safely removes all the loggers
	// All submanagers will also be closed
	Close() []error

	// NewSubManager creates a new manager that can have its own loggers.
	// The tag of the current manager will be passed to the submanager.
	// Calling a logging method on a submanager will trigger the same logging
	// method on the parent.
	NewSubManager(tag string) Manager

	// SetTag adds a tag to the logs
	SetTag(string)

	// Tag returns the tag of the manager
	Tag() string

	// FullTag returns the full tag (including parents) of the manager
	FullTag() string

	// Errorf logs an error message
	// Arguments are handled in the manner of fmt.Printf
	Errorf(msg string, args ...interface{})

	// Error logs an error message
	// Arguments are handled in the manner of fmt.Println.
	Error(args ...interface{})

	// Infof logs a message that may be helpful, but isn’t essential,
	// for troubleshooting
	// Arguments are handled in the manner of fmt.Printf
	Infof(msg string, args ...interface{})

	// Info logs a message that may be helpful, but isn’t essential,
	// for troubleshooting
	// Arguments are handled in the manner of fmt.Println.
	Info(args ...interface{})

	// Debugf logs a message that is intended for use in a development
	// environment while actively debugging your subsystem, not in shipping
	// software
	// Arguments are handled in the manner of fmt.Printf
	Debugf(msg string, args ...interface{})

	// Debug logs a message that is intended for use in a development
	// environment while actively debugging your subsystem, not in shipping
	// software
	// Arguments are handled in the manner of fmt.Println.
	Debug(args ...interface{})

	// Logf logs a message that might result a failure
	// Arguments are handled in the manner of fmt.Printf
	Logf(msg string, args ...interface{})

	// Log logs a message that might result a failure
	// Arguments are handled in the manner of fmt.Println.
	Log(args ...interface{})
}

// Err represents an error caused by a specific logger
type Err struct {
	error
	Logger Logger
}

// we make sure DefaultManager implements Manager
var _ Manager = (*DefaultManager)(nil)

// DefaultManager is a basic go-routine safe logger
type DefaultManager struct {
	sync.RWMutex

	id       string
	globals  map[string]interface{}
	loggers  map[string]Logger
	parent   *DefaultManager
	children map[string]*DefaultManager
	tag      string
}

// NewManager creates a new manager
func NewManager() Manager {
	return NewManagerWithTag("")
}

// NewManagerWithTag create a new manager tagged with the given tag
func NewManagerWithTag(tag string) Manager {
	return &DefaultManager{
		id:       uuid.New().String(),
		loggers:  map[string]Logger{},
		globals:  map[string]interface{}{},
		children: map[string]*DefaultManager{},
		tag:      tag,
	}
}

// AddGlobalData is used to add data that will be added to all logs
func (m *DefaultManager) AddGlobalData(key string, value interface{}) {
	m.Lock()
	defer m.Unlock()
	m.globals[key] = value
}

// RemoveGlobalData is used to remove data that are added to all logs
func (m *DefaultManager) RemoveGlobalData(key string) {
	m.Lock()
	defer m.Unlock()
	delete(m.globals, key)
}

// Add adds a logger
// returns ErrAlreadyExist if the logger has already been added
func (m *DefaultManager) Add(l Logger) error {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.loggers[l.ID()]; ok {
		return ErrAlreadyExist
	}

	m.loggers[l.ID()] = l
	return nil
}

// Remove safely removes a logger
// returns an Err struct if a logger could not be safely removed.
// Upon errors the logger will be force removed from the manager
func (m *DefaultManager) Remove(loggerID string) error {
	m.Lock()
	l, ok := m.loggers[loggerID]
	delete(m.loggers, loggerID)
	m.Unlock()

	if ok && l != nil {
		err := l.Close()
		if err != nil {
			return &Err{
				error:  err,
				Logger: l,
			}
		}
	}

	return nil
}

// Close safely removes all the loggers
// All submanagers will also be closed
// returns a list of errors if a logger could not be safely removed.
func (m *DefaultManager) Close() []error {
	return m.closeFromParent(false)
}

func (m *DefaultManager) closeFromParent(fromParents bool) []error {
	m.Lock()
	loggers := m.loggers
	m.loggers = map[string]Logger{}

	children := m.children
	m.children = map[string]*DefaultManager{}

	// if the parents is closing us, we don't need to ping it
	if !fromParents && m.parent != nil {
		m.parent.removeChild(m.ID())
	}
	m.Unlock()

	var errs []error
	// Close the loggers
	for _, l := range loggers {
		if err := l.Close(); err != nil {
			errs = append(errs, &Err{error: err, Logger: l})
		}
	}

	// We close the children too
	for _, c := range children {
		if cErrs := c.closeFromParent(true); cErrs != nil {
			errs = append(errs, cErrs...)
		}
	}

	return errs
}

// removeChild removes a child manager without closing it
func (m *DefaultManager) removeChild(id string) {
	m.Lock()
	delete(m.children, id)
	m.Unlock()
}

// NewSubManager creates a new manager that can have its own loggers.
// The tag of the current manager will be passed to the submanager.
// Calling a logging method on a submanager will trigger the same logging
// method on the parent.
func (m *DefaultManager) NewSubManager(tag string) Manager {
	m.Lock()
	defer m.Unlock()

	sm := NewManagerWithTag(tag)
	df := sm.(*DefaultManager)
	df.parent = m

	m.children[sm.ID()] = df
	return sm
}

// SetTag adds a tag to the logs
func (m *DefaultManager) SetTag(tag string) {
	m.Lock()
	defer m.Unlock()

	m.tag = tag
}

// Tag returns the tag of the manager
func (m *DefaultManager) Tag() string {
	m.RLock()
	defer m.RUnlock()
	return m.tag
}

// FullTag returns the full tag (including parents) of the manager
func (m *DefaultManager) FullTag() string {
	tag := m.Tag()
	if m.parent != nil {
		tag = m.parent.FullTag() + tag
	}
	return tag
}

// ID returns the manager's unique ID
func (m *DefaultManager) ID() string {
	// No need to lock since the ID should *never* be changed
	return m.id
}

// Errorf logs an error message
// Arguments are handled in the manner of fmt.Printf
func (m *DefaultManager) Errorf(msg string, args ...interface{}) {
	m.Error(fmt.Sprintf(msg, args...))
}

// Error logs an error message
// Arguments are handled in the manner of fmt.Println.
func (m *DefaultManager) Error(args ...interface{}) {
	msg := m.format(fmt.Sprintln(args...))
	m.error(msg)
}

func (m *DefaultManager) error(msg string) {
	m.RLock()
	defer m.RUnlock()

	// we send the log to the parent's logger first
	if m.parent != nil {
		m.parent.error(msg)
	}

	for _, l := range m.loggers {
		l.Error(msg)
	}
}

// Infof logs a message that may be helpful, but isn’t essential,
// for troubleshooting
// Arguments are handled in the manner of fmt.Printf
func (m *DefaultManager) Infof(msg string, args ...interface{}) {
	m.Info(fmt.Sprintf(msg, args...))
}

// Info logs a message that may be helpful, but isn’t essential,
// for troubleshooting
// Arguments are handled in the manner of fmt.Println.
func (m *DefaultManager) Info(args ...interface{}) {
	msg := m.format(fmt.Sprintln(args...))
	m.info(msg)
}

func (m *DefaultManager) info(msg string) {
	m.RLock()
	defer m.RUnlock()

	// we send the log to the parent's logger first
	if m.parent != nil {
		m.parent.info(msg)
	}

	for _, l := range m.loggers {
		l.Info(msg)
	}
}

// Debugf logs a message that is intended for use in a development
// environment while actively debugging your subsystem, not in shipping
// software
// Arguments are handled in the manner of fmt.Printf
func (m *DefaultManager) Debugf(msg string, args ...interface{}) {
	m.Debug(fmt.Sprintf(msg, args...))
}

// Debug logs a message that is intended for use in a development
// environment while actively debugging your subsystem, not in shipping
// software
// Arguments are handled in the manner of fmt.Println.
func (m *DefaultManager) Debug(args ...interface{}) {
	msg := m.format(fmt.Sprintln(args...))
	m.debug(msg)
}

func (m *DefaultManager) debug(msg string) {
	m.RLock()
	defer m.RUnlock()

	// we send the log to the parent's logger first
	if m.parent != nil {
		m.parent.debug(msg)
	}

	for _, l := range m.loggers {
		l.Debug(msg)
	}
}

// Logf logs a message that might result a failure
// Arguments are handled in the manner of fmt.Printf
func (m *DefaultManager) Logf(msg string, args ...interface{}) {
	m.Log(fmt.Sprintf(msg, args...))
}

// Log logs a message that might result a failure
// Arguments are handled in the manner of fmt.Println.
func (m *DefaultManager) Log(args ...interface{}) {
	msg := m.format(fmt.Sprintln(args...))
	m.log(msg)
}

func (m *DefaultManager) log(msg string) {
	m.RLock()
	defer m.RUnlock()

	// we send the log to the parent's logger first
	if m.parent != nil {
		m.parent.log(msg)
	}

	for _, l := range m.loggers {
		l.Log(msg)
	}
}

func (m *DefaultManager) format(msg string) string {
	tag := m.FullTag()
	if tag != "" {
		msg = tag + " " + msg
	}

	if msg[len(msg)-1] != '\n' {
		return msg + "\n"
	}

	globals := m.allGlobals()
	if len(globals) > 0 {
		jsonGlobals, err := json.Marshal(globals)
		if err != nil {
			panic(errors.Wrap(err, "could not encode the globals to JSON"))
		}
		msg += string(jsonGlobals) + "\n"
	}

	return msg
}

func (m *DefaultManager) allGlobals() map[string]interface{} {
	globals := map[string]interface{}{}
	if m.parent != nil {
		globals = m.parent.allGlobals()
	}
	for k, v := range m.globals {
		globals[k] = v
	}
	return globals
}
