# go-logger

[![Build Status](https://travis-ci.org/Nivl/go-logger.svg)](https://travis-ci.org/Nivl/go-logger)
[![Go Report Card](https://goreportcard.com/badge/github.com/nivl/go-logger)](https://goreportcard.com/report/github.com/nivl/go-logger)
[![GoDoc](https://godoc.org/github.com/Nivl/go-logger?status.svg)](https://godoc.org/github.com/Nivl/go-logger)

go-logger contains interfaces and basic implementations to deal with loggers

## Usage

```go
m := logger.NewManagerWithTag("[my-app]")

// Add a bunch of loggers
m.Add(logger.NewStderrLogger())
m.Add(NewFileLogger())

// send a log to all the loggers added with Add()
m.Errorf("error message: %s", "file not found") // prints "[ERROR][my-app] error message: file not found"

// create sub-loggers for specific parts of your app
// Sub-loggers can have their own loggers, but also reuse their parent's loggers
sm := m.NewSubLogger("[parser]")
sm.Log("foo") // prints "[my-app][parser] foo"
```

## Provided implementations

### gomock

```go
mockCtrl := gomock.NewController(t)
defer mockCtrl.Finish()

m := mocklogger.NewMockManager(mockCtrl)
m.EXPECTS().Log("foo")

```

### StderrLogger (log.Print() wrapper)

```go
m := logger.NewManager()
m.Add(logger.NewStderrLogger())

m.Error("foobar") // printed on stderr
```

### External implementations

- [Native Logger](https://github.com/Nivl/gologger-native): Logger using the native log system of the current OS
